package brucecore

import (
	"strings"

	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// HostnameInfo - hostname infomation
type HostnameInfo struct {
	Hostname string
	IPAddr   []string
}

// Insert - insert a SourceInfo
func (hi *HostnameInfo) Insert(ipaddr string) {
	if len(hi.IPAddr) > 0 {
		for _, cip := range hi.IPAddr {
			if cip == ipaddr {
				return
			}
		}
	}

	hi.IPAddr = append(hi.IPAddr, ipaddr)
}

// HostnameList - hostname infomation list
type HostnameList struct {
	List []*HostnameInfo
}

// Insert - insert a SourceInfo
func (sl *HostnameList) Insert(hostname string, ipaddr []string) {
	v := sl.Find(hostname)
	if v == nil {
		sl.List = append(sl.List, &HostnameInfo{
			Hostname: hostname,
			IPAddr:   ipaddr,
		})

		return
	}

	for _, curip := range ipaddr {
		v.Insert(curip)
	}
}

// Find - find a SourceInfo
func (sl *HostnameList) Find(hostname string) *HostnameInfo {
	for _, v := range sl.List {
		if v.Hostname == hostname {
			return v
		}
	}

	return nil
}

// ToData - to data
func (sl *HostnameList) ToData() [][]string {
	var lst [][]string

	for _, v := range sl.List {
		lst = append(lst, []string{
			v.Hostname,
			strings.Join(v.IPAddr, ";"),
		})
	}

	return lst
}

// AnalyzeHostNameInfo - analyze hostname
func AnalyzeHostNameInfo(reqs []*jarviscrawlercore.AnalyzeReqInfo) (*HostnameList, error) {
	lst := &HostnameList{}

	for _, v := range reqs {
		cs, err := GetHostname(v.Url)
		if err != nil {
			return nil, err
		}

		ips := strings.Split(v.Ipaddr, ";")
		lst.Insert(cs, ips)
	}

	return lst, nil
}
