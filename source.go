package brucecore

import jarviscrawlercore "github.com/zhs007/jccclient/proto"

// SourceInfo - source infomation
type SourceInfo struct {
	Source     string
	TotalBytes int
	TotalTime  int
	TotalNums  int
}

// SourceList - source infomation list
type SourceList struct {
	List []*SourceInfo
}

// ResSourceData - resource source data
type ResSourceData struct {
	ResSource     []string  `yaml:"source"`
	Bytes         []int     `yaml:"bytes"`
	MBytes        []float32 `yaml:"mbytes"`
	Time          []int     `yaml:"time"`
	Nums          []int     `yaml:"nums"`
	DownloadSpeed []float32 `yaml:"downloadspeed"`
}

// Insert - insert a SourceInfo
func (sl *SourceList) Insert(host string, bytes int, time int) {
	v := sl.Find(host)
	if v == nil {
		sl.List = append(sl.List, &SourceInfo{
			Source:     host,
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
func (sl *SourceList) Find(host string) *SourceInfo {
	for _, v := range sl.List {
		if v.Source == host {
			return v
		}
	}

	return nil
}

// ToData - SourceList => ResSourceData
func (sl *SourceList) ToData() *ResSourceData {
	rsd := &ResSourceData{}

	for _, v := range sl.List {
		rsd.ResSource = append(rsd.ResSource, v.Source)
		rsd.Bytes = append(rsd.Bytes, v.TotalBytes)
		rsd.MBytes = append(rsd.MBytes, float32(v.TotalBytes)/1024/1024)
		rsd.Time = append(rsd.Time, v.TotalTime)
		rsd.Nums = append(rsd.Nums, v.TotalNums)

		ds := float32(v.TotalBytes) * 1000 / float32(v.TotalTime) / 1024 / 1024

		rsd.DownloadSpeed = append(rsd.DownloadSpeed, ds)
	}

	return rsd
}

// AnalyzeResSource - analyze request
func AnalyzeResSource(reqs []*jarviscrawlercore.AnalyzeReqInfo) (*SourceList, error) {
	sl := &SourceList{}

	for _, v := range reqs {
		cs, err := GetHostname(v.Url)
		if err != nil {
			return nil, err
		}

		sl.Insert(cs, int(v.BufBytes), int(v.DownloadTime))
	}

	return sl, nil
}
