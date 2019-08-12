package main

import (
	"context"
	"fmt"

	"github.com/zhs007/adacore"
	adacorepb "github.com/zhs007/adacore/proto"
	"github.com/zhs007/brucecore"
	"github.com/zhs007/brucecore/templates/spnormal"
	"github.com/zhs007/jccclient"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

func analyzePage(url string, delay int, w int, h int) (*jarviscrawlercore.ReplyAnalyzePage, error) {
	// client := jccclient.NewClient("127.0.0.1:7051", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")
	client := jccclient.NewClient("47.75.11.61:7051", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")

	reply, err := client.AnalyzePage(context.Background(),
		url, delay, w, h)

	if err != nil {
		// fmt.Printf("AnalyzePage %v", err)

		return nil, err
	}

	// if reply != nil {
	// 	fmt.Printf("%v", reply)
	// }

	return reply, nil
}

func genMarkdown(url string, reply *jarviscrawlercore.ReplyAnalyzePage) *adacorepb.MarkdownData {
	mddata := &adacorepb.MarkdownData{
		TemplateName: "default",
	}

	km, err := adacore.LoadKeywordMappingList("./keywordmapping.yaml")
	if err != nil {
		fmt.Printf("load keywordmapping error %v", err)
	}

	md := adacore.NewMakrdown("Analyze Page Result")

	md.AppendTable([]string{"Title", "Infomation"}, [][]string{
		[]string{"URL", "[click here](http://47.90.46.159:8090/game.html?gameCode=nightclub&language=zh_CN&isCheat=true&slotKey=)"},
		[]string{"Loading Time", brucecore.FormatTime(int(reply.PageTime))},
		[]string{"Resource Nums", fmt.Sprintf("%v", len(reply.Reqs))},
		[]string{"Total Resource Size", brucecore.FormatByteSize(int(reply.PageBytes))},
	})

	md.AppendParagraph("This libraray is write by Zerro.\nThis is a multi-line text.")

	for i, v := range reply.Screenshots {
		md.AppendImageBuf("Screenshot", fmt.Sprintf("screenshot%v", i), v.Buf, mddata)
	}

	mddata.StrData = md.GetMarkdownString(km)

	// fmt.Print(str)

	return mddata
}

func requestAda(url string, result *jarviscrawlercore.ReplyAnalyzePage) error {
	client := adacore.NewClient("47.91.209.141:7201", "x7sSGGHgmKwUMoa5S4VZlr9bUF2lCCzF")

	km, err := adacore.LoadKeywordMappingList("./keywordmapping.yaml")
	if err != nil {
		fmt.Printf("load keywordmapping error %v", err)

		return err
	}

	mddata, err := spnormal.GenMarkdown("夜店", url, result, km)
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

	adacore.InitTemplates(cfg.TemplatesPath)
	adacore.InitLogger(cfg)

	reply, err := analyzePage("http://47.90.46.159:8090/game.html?gameCode=nightclub&language=zh_CN&isCheat=true&slotKey=",
		10, 1280, 800)
	if err != nil {
		fmt.Printf("analyzePage err %v", err)
	}

	if reply != nil {
		fmt.Printf("analyzePage ok!\n")

		err = requestAda("http://47.90.46.159:8090/game.html?gameCode=nightclub&language=zh_CN&isCheat=true&slotKey=", reply)
		if err != nil {
			fmt.Printf("requestAda err %v", err)
		}
	}
}
