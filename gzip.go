package brucecore

import (
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// GetNoGZip - get no gzip resource
func GetNoGZip(reqs []*jarviscrawlercore.AnalyzeReqInfo) (
	[]string, error) {

	var lst []string

	for _, v := range reqs {
		rt, err := GetResType(v)
		if err != nil {
			return nil, err
		}

		if rt == "mp3" || rt == "webp" || rt == "gif" ||
			rt == "jpg" || rt == "png" || rt == "mp4" {

			continue
		}

		if v.ImgWidth == 0 && v.ImgHeight == 0 {
			if !v.IsGZip {
				lst = append(lst, v.Url)
			}
		}
	}

	return lst, nil
}
