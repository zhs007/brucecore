package brucecore

import (
	"fmt"
	"strings"

	"github.com/zhs007/adacore"

	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// GetURLName - get URL name
func GetURLName(url string) string {
	arr := strings.Split(strings.Split(url, "?")[0], "/")

	return arr[len(arr)-1]
}

type reqImgType struct {
	w   int
	h   int
	lst []*jarviscrawlercore.AnalyzeReqInfo
}

// buildData - build data
func (rit *reqImgType) buildData(data *adacore.ChartTreeMapData) {
	for _, v := range rit.lst {
		curd := adacore.ChartTreeMapData{
			Name:  GetURLName(v.Url),
			Value: int(v.BufBytes),
			URL:   v.Url,
		}

		data.Children = append(data.Children, curd)
	}
}

type reqInType struct {
	rtype  string
	lst    []*jarviscrawlercore.AnalyzeReqInfo
	lstImg []*reqImgType
}

// insert - insert a reqImgType
func (rit *reqInType) insert(w int, h int, req *jarviscrawlercore.AnalyzeReqInfo) {
	w1 := w
	h1 := h
	if h > w {
		w1 = h
		h1 = w
	}

	v := rit.find(w1, h1)
	if v == nil {
		rit.lstImg = append(rit.lstImg, &reqImgType{
			w: w1,
			h: h1,
			lst: []*jarviscrawlercore.AnalyzeReqInfo{
				req,
			},
		})

		return
	}

	v.lst = append(v.lst, req)
}

// find - find a reqImgType
func (rit *reqInType) find(w int, h int) *reqImgType {
	for _, v := range rit.lstImg {
		if v.w == w && v.h == h {
			return v
		}
	}

	return nil
}

// buildImg - build image
func (rit *reqInType) buildImg() {
	for _, v := range rit.lst {
		if v.ImgWidth > 0 && v.ImgHeight > 0 {
			rit.insert(int(v.ImgWidth), int(v.ImgHeight), v)
		}
	}
}

// buildData - build data
func (rit *reqInType) buildData(data *adacore.ChartTreeMapData) {
	if len(rit.lstImg) > 0 {
		for _, v := range rit.lstImg {
			curd := adacore.ChartTreeMapData{
				Name: fmt.Sprintf("%vx%v", v.w, v.h),
			}

			v.buildData(&curd)

			data.Children = append(data.Children, curd)
		}

		return
	}

	for _, v := range rit.lst {
		curd := adacore.ChartTreeMapData{
			Name:  GetURLName(v.Url),
			Value: int(v.BufBytes),
			URL:   v.Url,
		}

		data.Children = append(data.Children, curd)
	}
}

type reqInSource struct {
	source  string
	lst     []*jarviscrawlercore.AnalyzeReqInfo
	lstType []*reqInType
}

// insert - insert a reqInType
func (ris *reqInSource) insert(rtype string, req *jarviscrawlercore.AnalyzeReqInfo) {
	v := ris.find(rtype)
	if v == nil {
		ris.lstType = append(ris.lstType, &reqInType{
			rtype: rtype,
			lst: []*jarviscrawlercore.AnalyzeReqInfo{
				req,
			},
		})

		return
	}

	v.lst = append(v.lst, req)
}

// find - find a reqInType
func (ris *reqInSource) find(rtype string) *reqInType {
	for _, v := range ris.lstType {
		if v.rtype == rtype {
			return v
		}
	}

	return nil
}

// buildResType - build resource type
func (ris *reqInSource) buildResType() {
	for _, v := range ris.lst {
		rt, err := GetResType(v)
		if err == nil {
			ris.insert(rt, v)
		}
	}

	for _, v := range ris.lstType {
		v.buildImg()
	}
}

// buildData - build data
func (ris *reqInSource) buildData(data *adacore.ChartTreeMapSeriesNode) {
	if ris.source == "local:imgdata" {
		// for _, v := range ris.lst {
		// 	data.Value += int(v.BufBytes)
		// }

		return
	}

	for _, v := range ris.lstType {
		curd := adacore.ChartTreeMapData{
			Name: v.rtype,
		}

		v.buildData(&curd)

		data.Data = append(data.Data, curd)
	}
}

// ResTreeMapData - resource treemap data
type ResTreeMapData struct {
	TreeMap  []adacore.ChartTreeMapSeriesNode
	HostList []string

	lstReqInSource []*reqInSource
}

// insert - insert a SourceInfo
func (data *ResTreeMapData) insert(host string, req *jarviscrawlercore.AnalyzeReqInfo) {
	v := data.find(host)
	if v == nil {
		data.lstReqInSource = append(data.lstReqInSource, &reqInSource{
			source: host,
			lst: []*jarviscrawlercore.AnalyzeReqInfo{
				req,
			},
		})

		return
	}

	v.lst = append(v.lst, req)
}

// find - find a SourceInfo
func (data *ResTreeMapData) find(host string) *reqInSource {
	for _, v := range data.lstReqInSource {
		if v.source == host {
			return v
		}
	}

	return nil
}

// buildData - build data
func (data *ResTreeMapData) buildData() {
	for _, v := range data.lstReqInSource {
		if v.source == "local:imgdata" {
			continue
		}

		curd := adacore.ChartTreeMapSeriesNode{
			Name: v.source,
		}

		v.buildData(&curd)

		data.TreeMap = append(data.TreeMap, curd)
		data.HostList = append(data.HostList, v.source)
	}
}

// AnalyzeResTreeMap - analyze request
func AnalyzeResTreeMap(reqs []*jarviscrawlercore.AnalyzeReqInfo) (*ResTreeMapData, error) {
	data := &ResTreeMapData{}

	for _, v := range reqs {
		cs, err := GetHostname(v.Url)
		if err != nil {
			return nil, err
		}

		data.insert(cs, v)
	}

	for _, v := range data.lstReqInSource {
		v.buildResType()
	}

	data.buildData()

	return data, nil
}
