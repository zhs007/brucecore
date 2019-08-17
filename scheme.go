package brucecore

import (
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// GetResWithScheme - get resource with scheme
func GetResWithScheme(reqs []*jarviscrawlercore.AnalyzeReqInfo, scheme string) (
	[]string, error) {

	var lst []string

	for _, v := range reqs {
		cs, _ := GetScheme(v.Url)
		if cs == scheme {
			lst = append(lst, v.Url)
		}
	}

	return lst, nil
}
