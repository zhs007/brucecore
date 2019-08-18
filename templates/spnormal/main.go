package spnormal

import (
	"fmt"
	"strings"

	"github.com/zhs007/adacore"
	adacorepb "github.com/zhs007/adacore/proto"
	"github.com/zhs007/brucecore"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// GenMarkdown - generate markdown
func GenMarkdown(name string, url string, reply *jarviscrawlercore.ReplyAnalyzePage,
	km *adacore.KeywordMappingList) (*adacorepb.MarkdownData, error) {

	mddata := &adacorepb.MarkdownData{
		TemplateName: "default",
	}

	sl, err := brucecore.AnalyzeResSource(reply.Reqs)
	if err != nil {
		return nil, err
	}

	rt, err := brucecore.AnalyzeResType(reply.Reqs)
	if err != nil {
		return nil, err
	}

	restreemap, err := brucecore.AnalyzeResTreeMap(reply.Reqs)
	if err != nil {
		return nil, err
	}

	httpcheme, err := brucecore.GetResWithScheme(reply.Reqs, "http")
	if err != nil {
		return nil, err
	}

	nogzip, err := brucecore.GetNoGZip(reply.Reqs)
	if err != nil {
		return nil, err
	}

	hostname, err := brucecore.AnalyzeHostNameInfo(reply.Reqs)
	if err != nil {
		return nil, err
	}

	md := adacore.NewMakrdown(fmt.Sprintf("单页分析结果 - %v", name))

	md.AppendTable([]string{"常规项目", "结果"}, [][]string{
		[]string{"测试地址", fmt.Sprintf("[%v](%v)", name, url)},
		[]string{"载入总时长", brucecore.FormatTime(int(reply.PageTime))},
		[]string{"加载资源数量", fmt.Sprintf("%v", len(reply.Reqs))},
		[]string{"加载资源总大小", brucecore.FormatByteSize(int(reply.PageBytes))},
		[]string{"平均加载速度", brucecore.FormatByteSize(int(int64(reply.PageBytes)*1000/int64(reply.PageTime))) + "/s"},
		[]string{"错误数量", fmt.Sprintf("%v", len(reply.Errs))},
		[]string{"日志数量", fmt.Sprintf("%v", len(reply.Logs))},
	})

	md.AppendParagraph("")

	md.AppendTable([]string{"hostname", "IP"}, hostname.ToData())

	md.AppendParagraph("### 日志")
	md.AppendCode("", strings.Join(reply.Logs, "\n"))

	md.AppendParagraph("### 错误输出")
	md.AppendCode("", strings.Join(reply.Errs, "\n"))

	if len(httpcheme) > 0 {
		md.AppendParagraph("### http协议")
		md.AppendCode("", strings.Join(httpcheme, "\n"))
	}

	if len(nogzip) > 0 {
		md.AppendParagraph("### 未GZip")

		md.AppendTable([]string{"URL", "原文件大小", "压缩后大小", "降低比例"},
			brucecore.BuildNoGZipTableData(nogzip))
		// md.AppendCode("", strings.Join(nogzip, "\n"))
	}

	_, err = md.AppendDataset("reshostds", sl.ToData())
	if err != nil {
		return nil, err
	}

	_, err = md.AppendChartBar(&adacore.ChartBar{
		ID:          "reshostspeed",
		DatasetName: "reshostds",
		Title:       "下载来源站速度比较",
		SubText:     "",
		Width:       1280,
		Height:      800,
		LegendData:  []string{"下载速度", "下载量"},
		XType:       "category",
		XData:       "source",
		XShowAll:    true,
		YType:       "value",
		YData: []adacore.ChartBasicData{
			adacore.ChartBasicData{
				Name: "下载速度",
				Data: "downloadspeed",
			},
			adacore.ChartBasicData{
				Name: "下载量",
				Data: "mbytes",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "reshostbytes",
	// 	DatasetName: "reshostds",
	// 	Title:       "下载资源大小分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "下载来源",
	// 	BVal:        "source",
	// 	CVal:        "bytes",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "reshostnums",
	// 	DatasetName: "reshostds",
	// 	Title:       "下载资源数量分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "下载来源",
	// 	BVal:        "source",
	// 	CVal:        "nums",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "reshosttime",
	// 	DatasetName: "reshostds",
	// 	Title:       "下载资源耗时分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "下载来源",
	// 	BVal:        "source",
	// 	CVal:        "time",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	_, err = md.AppendDataset("rtds", rt.ToData())
	if err != nil {
		return nil, err
	}

	_, err = md.AppendChartBar(&adacore.ChartBar{
		ID:          "rtspeed",
		DatasetName: "rtds",
		Title:       "资源类型下载速度比较",
		SubText:     "",
		Width:       1280,
		Height:      800,
		LegendData:  []string{"下载速度", "下载量"},
		XType:       "category",
		XData:       "restype",
		XShowAll:    true,
		YType:       "value",
		YData: []adacore.ChartBasicData{
			adacore.ChartBasicData{
				Name: "下载速度",
				Data: "downloadspeed",
			},
			adacore.ChartBasicData{
				Name: "下载量",
				Data: "mbytes",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "rtbytes",
	// 	DatasetName: "rtds",
	// 	Title:       "资源类型大小分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "资源类型",
	// 	BVal:        "restype",
	// 	CVal:        "bytes",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "rtnums",
	// 	DatasetName: "rtds",
	// 	Title:       "资源类型数量分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "下载来源",
	// 	BVal:        "restype",
	// 	CVal:        "nums",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "rttime",
	// 	DatasetName: "rtds",
	// 	Title:       "资源类型耗时分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "下载来源",
	// 	BVal:        "restype",
	// 	CVal:        "time",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "imgbytes",
	// 	DatasetName: "rtds",
	// 	Title:       "图片大小分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "图片规格",
	// 	BVal:        "imgtype",
	// 	CVal:        "imgbytes",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "imgnums",
	// 	DatasetName: "rtds",
	// 	Title:       "图片数量分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "图片规格",
	// 	BVal:        "imgtype",
	// 	CVal:        "imgnums",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = md.AppendChartPie(&adacore.ChartPie{
	// 	ID:          "imgtime",
	// 	DatasetName: "rtds",
	// 	Title:       "图片耗时分布",
	// 	SubText:     "",
	// 	Width:       1280,
	// 	Height:      800,
	// 	A:           "图片规格",
	// 	BVal:        "imgtype",
	// 	CVal:        "imgtime",
	// })
	// if err != nil {
	// 	return nil, err
	// }

	_, err = md.AppendChartTreeMap(&adacore.ChartTreeMap{
		ID:         "restreemap",
		Title:      "资源大小分布",
		SubText:    "",
		Width:      800,
		Height:     600,
		LegendData: restreemap.HostList,
		TreeMap:    restreemap.TreeMap,
	})
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%v", str1)

	for i, v := range reply.Screenshots {
		imgname := fmt.Sprintf("screenshot%v.jpg", i)
		if v.Type == jarviscrawlercore.AnalyzeScreenshotType_AST_PNG {
			imgname = fmt.Sprintf("screenshot%v.png", i)
		}

		md.AppendImageBuf("截图", imgname, v.Buf, mddata)
	}

	mddata.StrData = md.GetMarkdownString(km)

	return mddata, nil
}
