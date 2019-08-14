package brucecore

import (
	"fmt"
	"strings"

	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// ResTypeInfo - resource type infomation
type ResTypeInfo struct {
	ResType    string
	TotalBytes int
	TotalTime  int
	TotalNums  int
}

// ImgTypeInfo - image type infomation
type ImgTypeInfo struct {
	ImgWidth   int
	ImgHeight  int
	TotalBytes int
	TotalTime  int
	TotalNums  int
}

// ResTypeMgr - resource type infomation list
type ResTypeMgr struct {
	List    []*ResTypeInfo
	ImgList []*ImgTypeInfo
}

// ResTypeData - resource type data
type ResTypeData struct {
	ResType       []string  `yaml:"restype"`
	Bytes         []int     `yaml:"bytes"`
	MBytes        []float32 `yaml:"mbytes"`
	Time          []int     `yaml:"time"`
	Nums          []int     `yaml:"nums"`
	DownloadSpeed []float32 `yaml:"downloadspeed"`
	ImgType       []string  `yaml:"imgtype"`
	ImgBytes      []int     `yaml:"imgbytes"`
	ImgTime       []int     `yaml:"imgtime"`
	ImgNums       []int     `yaml:"imgnums"`
}

// Insert - insert a SourceInfo
func (mgr *ResTypeMgr) Insert(restype string, bytes int, time int) {
	v := mgr.Find(restype)
	if v == nil {
		mgr.List = append(mgr.List, &ResTypeInfo{
			ResType:    restype,
			TotalBytes: bytes,
			TotalTime:  time,
			TotalNums:  1,
		})

		return
	}

	v.TotalBytes += bytes
	v.TotalTime += time
	v.TotalNums++
}

// Find - find a SourceInfo
func (mgr *ResTypeMgr) Find(restype string) *ResTypeInfo {
	for _, v := range mgr.List {
		if v.ResType == restype {
			return v
		}
	}

	return nil
}

// InsertImage - insert a image
func (mgr *ResTypeMgr) InsertImage(w int, h int, bytes int, time int) {
	v := mgr.FindImage(w, h)
	if v == nil {
		mgr.ImgList = append(mgr.ImgList, &ImgTypeInfo{
			ImgWidth:   w,
			ImgHeight:  h,
			TotalBytes: bytes,
			TotalTime:  time,
			TotalNums:  1,
		})

		return
	}

	v.TotalBytes += bytes
	v.TotalTime += time
	v.TotalNums++
}

// FindImage - find a image
func (mgr *ResTypeMgr) FindImage(w int, h int) *ImgTypeInfo {
	for _, v := range mgr.ImgList {
		if v.ImgWidth == w && v.ImgHeight == h {
			return v
		}
	}

	return nil
}

// ToData - SourceList => ResSourceData
func (mgr *ResTypeMgr) ToData() *ResTypeData {
	rtd := &ResTypeData{}

	for _, v := range mgr.List {
		rtd.ResType = append(rtd.ResType, v.ResType)
		rtd.Bytes = append(rtd.Bytes, v.TotalBytes)
		rtd.MBytes = append(rtd.MBytes, float32(v.TotalBytes)/1024/1024)
		rtd.Time = append(rtd.Time, v.TotalTime)
		rtd.Nums = append(rtd.Nums, v.TotalNums)

		ds := float32(v.TotalBytes) * 1000 / float32(v.TotalTime) / 1024 / 1024

		rtd.DownloadSpeed = append(rtd.DownloadSpeed, ds)
	}

	for _, v := range mgr.ImgList {
		rtd.ImgType = append(rtd.ImgType, fmt.Sprintf("%vx%v", v.ImgWidth, v.ImgHeight))
		rtd.ImgBytes = append(rtd.ImgBytes, v.TotalBytes)
		rtd.ImgTime = append(rtd.ImgTime, v.TotalTime)
		rtd.ImgNums = append(rtd.ImgNums, v.TotalNums)
	}

	return rtd
}

// AnalyzeResType - analyze request
func AnalyzeResType(reqs []*jarviscrawlercore.AnalyzeReqInfo) (*ResTypeMgr, error) {
	mgr := &ResTypeMgr{}

	for _, v := range reqs {
		cct := strings.Split(v.ContentType, ";")
		cft := strings.Split(cct[0], "/")
		if cft[1] == "jpg" || cft[1] == "jpeg" {
			cft[1] = "jpg"
		} else if cft[1] == "javascript" {
			cft[1] = "js"
		}

		mgr.Insert(cft[1], int(v.BufBytes), int(v.DownloadTime))

		if v.ImgWidth > 0 && v.ImgHeight > 0 {
			mgr.InsertImage(int(v.ImgWidth), int(v.ImgHeight), int(v.BufBytes), int(v.DownloadTime))
		}
	}

	return mgr, nil
}
