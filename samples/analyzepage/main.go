package main

import (
	"context"
	"fmt"

	"github.com/zhs007/jccclient"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

func analyzePage(url string, delay int, w int, h int) (*jarviscrawlercore.ReplyAnalyzePage, error) {
	client := jccclient.NewClient("127.0.0.1:7051", "wzDkh9h2fhfUVuS9jZ8uVbhV3vC5AWX3")

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

func main() {
	reply, err := analyzePage("http://47.90.46.159:8090/game.html?gameCode=nightclub&language=zh_CN&isCheat=true&slotKey=",
		10, 1280, 800)
	if err != nil {
		fmt.Printf("analyzePage err %v", err)
	}

	if reply != nil {
		fmt.Printf("analyzePage ok!")
	}
}
