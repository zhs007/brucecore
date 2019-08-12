package spnormal

import (
	"fmt"

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

	md := adacore.NewMakrdown(fmt.Sprintf("单页分析结果 - %v", name))

	md.AppendTable([]string{"常规项目", "结果"}, [][]string{
		[]string{"测试地址", fmt.Sprintf("[%v](%v)", name, url)},
		[]string{"载入总时长", brucecore.FormatTime(int(reply.PageTime))},
		[]string{"加载资源数量", fmt.Sprintf("%v", len(reply.Reqs))},
		[]string{"加载资源总大小", brucecore.FormatByteSize(int(reply.PageBytes))},
		[]string{"平均加载速度", brucecore.FormatByteSize(int(int64(reply.PageBytes)*1000/int64(reply.PageTime))) + "/s"},
	})

	sl, err := brucecore.AnalyzeResSource(reply.Reqs)
	if err != nil {
		return nil, err
	}

	_, err = md.AppendDataset("reshostds", sl.ToData())
	if err != nil {
		return nil, err
	}

	_, err = md.AppendChartPie(&adacore.ChartPie{
		ID:          "reshostbytes",
		DatasetName: "reshostds",
		Title:       "下载资源大小分布",
		SubText:     "",
		Width:       1280,
		Height:      800,
		A:           "下载来源",
		BVal:        "source",
		CVal:        "bytes",
	})
	if err != nil {
		return nil, err
	}

	_, err = md.AppendChartPie(&adacore.ChartPie{
		ID:          "reshostnums",
		DatasetName: "reshostds",
		Title:       "下载资源数量分布",
		SubText:     "",
		Width:       1280,
		Height:      800,
		A:           "下载来源",
		BVal:        "source",
		CVal:        "nums",
	})
	if err != nil {
		return nil, err
	}

	_, err = md.AppendChartPie(&adacore.ChartPie{
		ID:          "reshosttime",
		DatasetName: "reshostds",
		Title:       "下载资源耗时分布",
		SubText:     "",
		Width:       1280,
		Height:      800,
		A:           "下载来源",
		BVal:        "source",
		CVal:        "time",
	})
	if err != nil {
		return nil, err
	}

	for i, v := range reply.Screenshots {
		md.AppendImageBuf("截图", fmt.Sprintf("screenshot%v", i), v.Buf, mddata)
	}

	mddata.StrData = md.GetMarkdownString(km)

	return mddata, nil
}
