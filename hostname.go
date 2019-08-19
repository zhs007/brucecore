package brucecore

import (
	"context"
	"fmt"
	"strings"

	"github.com/zhs007/brucecore/ipgeo"
	ipgeopb "github.com/zhs007/brucecore/ipgeo/proto"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// IPGeoInfo - ipgeo infomation
type IPGeoInfo struct {
	IPAddr string
	IPGeo  *ipgeopb.IPGeo
}

// HostnameInfo - hostname infomation
type HostnameInfo struct {
	Hostname string
	IPs      []*IPGeoInfo
}

// Insert - insert a SourceInfo
func (hi *HostnameInfo) Insert(ipaddr string, ipg *ipgeopb.IPGeo) {
	if len(hi.IPs) > 0 {
		for _, cip := range hi.IPs {
			if cip.IPAddr == ipaddr {
				return
			}
		}
	}

	hi.IPs = append(hi.IPs, &IPGeoInfo{
		IPAddr: ipaddr,
		IPGeo:  ipg,
	})
}

// HostnameList - hostname infomation list
type HostnameList struct {
	List []*HostnameInfo
}

// Insert - insert a SourceInfo
func (sl *HostnameList) Insert(ctx context.Context, ipgeodb *ipgeo.DB, hostname string, ipaddr []string) {
	v := sl.Find(hostname)
	if v == nil {
		sl.List = append(sl.List, &HostnameInfo{
			Hostname: hostname,
		})

		v = sl.Find(hostname)
		if v == nil {
			return
		}
	}

	for _, curip := range ipaddr {
		ipg, err := ipgeodb.GetIPGeoEx(ctx, curip)
		if err == nil && ipg != nil {
			v.Insert(curip, ipg)
		}
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
		ips := []string{}
		for _, ipv := range v.IPs {
			ips = append(ips, fmt.Sprintf("%v %v.%v", ipv.IPAddr, ipv.IPGeo.Continent, ipv.IPGeo.Country))
		}

		lst = append(lst, []string{
			v.Hostname,
			strings.Join(ips, ";"),
		})
	}

	return lst
}

// AnalyzeHostNameInfo - analyze hostname
func AnalyzeHostNameInfo(ctx context.Context, ipgeodb *ipgeo.DB,
	reqs []*jarviscrawlercore.AnalyzeReqInfo) (*HostnameList, error) {

	lst := &HostnameList{}

	for _, v := range reqs {
		cs, err := GetHostname(v.Url)
		if err != nil {
			return nil, err
		}

		if cs != "" {
			ips := strings.Split(v.Ipaddr, ";")
			lst.Insert(ctx, ipgeodb, cs, ips)
		}
	}

	return lst, nil
}
