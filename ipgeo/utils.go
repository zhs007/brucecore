package ipgeo

import (
	brucecorebase "github.com/zhs007/brucecore/base"
)

// IPGeoDBKeyPrefix - This is the prefix of IPGeoDBKey
const IPGeoDBKeyPrefix = "ip:"

// makeKey - Generate a database key via ip
func makeKey(ip string) string {
	return brucecorebase.AppendString(IPGeoDBKeyPrefix, ip)
}
