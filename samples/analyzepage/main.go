package main

import (
	"context"
	"fmt"

	"github.com/zhs007/adacore"
	brucetemplates "github.com/zhs007/brucecore/templates"
	"github.com/zhs007/jccclient"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

func analyzePage(url string, delay int, w int, h int) (*jarviscrawlercore.ReplyAnalyzePage, error) {
	// client := jccclient.NewClient("127.0.0.1:7051", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")
	client := jccclient.NewClient("47.75.11.61:7051", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")

	reply, err := client.AnalyzePage(context.Background(),
		url, &jccclient.Viewport{
			Width:             w,
			Height:            h,
			DeviceScaleFactor: 1.0,
			IsMobile:          false,
			IsLandscape:       false,
		},
		&jccclient.AnalyzePageOptions{
			NeedScreenshots:  true,
			NeedLogs:         true,
			Timeout:          0,
			ScreenshotsDelay: delay,
		})

	if err != nil {
		// fmt.Printf("AnalyzePage %v", err)

		return nil, err
	}

	// if reply != nil {
	// 	fmt.Printf("%v", reply)
	// }

	return reply, nil
}

func requestAda(name string, url string, result *jarviscrawlercore.ReplyAnalyzePage) error {
	client := adacore.NewClient("47.91.209.141:7201", "x7sSGGHgmKwUMoa5S4VZlr9bUF2lCCzF")

	km, err := adacore.LoadKeywordMappingList("./keywordmapping.yaml")
	if err != nil {
		fmt.Printf("load keywordmapping error %v", err)

		return err
	}

	mddata, err := brucetemplates.GenMarkdown("spnormal", name, url, result, km)
	if err != nil {
		fmt.Printf("spnormal.GenMarkdown error %v", err)

		return err
	}

	reply, err := client.BuildWithMarkdown(context.Background(), mddata)
	if err != nil {
		fmt.Printf("startClient BuildWithMarkdownFile %v", err)

		return err
	}

	if reply != nil {
		// fmt.Print(reply.HashName)
		fmt.Print(reply.Url)
	}

	return nil
}

func main() {
	cfg, err := adacore.LoadConfig("./adacore.yaml")
	if err != nil {
		fmt.Printf("startServ LoadConfig %v", err)

		return
	}

	url := "https://www.douban.com"
	name := "豆瓣"

	adacore.InitTemplates()
	adacore.InitLogger(cfg)

	reply, err := analyzePage(url, 10, 1280, 800)
	if err != nil {
		fmt.Printf("analyzePage err %v", err)
	}

	if reply != nil {
		fmt.Printf("analyzePage ok!\n")

		err = requestAda(name, url, reply)
		if err != nil {
			fmt.Printf("requestAda err %v", err)
		}
	}
}
