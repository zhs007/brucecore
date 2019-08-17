package brucecore

import (
	"fmt"
	"strings"

	"net/url"

	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// FormatTime - 12345 => 1.235s
func FormatTime(time int) string {
	if time > 24*60*60*1000 {
		d := float64(time) / (24 * 60 * 60 * 1000)

		return fmt.Sprintf("%.3f D", d)
	} else if time > 60*60*1000 {
		h := float64(time) / (60 * 60)

		return fmt.Sprintf("%.3f H", h)
	} else if time > 60*1000 {
		m := float64(time) / 60

		return fmt.Sprintf("%.3f m", m)
	}

	return fmt.Sprintf("%.3f s", float64(time)/1000)
}

// FormatByteSize - 1025 => 1k1b
func FormatByteSize(bytesize int) string {
	if bytesize > 1024*1024*1024 {
		g := float64(bytesize) / (1024 * 1024 * 1024)

		return fmt.Sprintf("%.3f G", g)
	} else if bytesize > 1024*1024 {
		m := float64(bytesize) / (1024 * 1024)

		return fmt.Sprintf("%.3f M", m)
	} else if bytesize > 1024 {
		k := float64(bytesize) / 1024

		return fmt.Sprintf("%.3f K", k)
	}

	return fmt.Sprintf("%v B", bytesize)
}

// GetHostname - https://www.a.com/b/c.png => www.a.com
func GetHostname(str string) (string, error) {
	if strings.Index(str, "local:") == 0 {
		return "local", nil
	}

	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}

// GetScheme - https://www.a.com/b/c.png => https
func GetScheme(str string) (string, error) {
	if strings.Index(str, "local:") == 0 {
		return "", nil
	}

	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	return u.Scheme, nil
}

// GetResType - https://www.a.com/b/c.png => https
func GetResType(res *jarviscrawlercore.AnalyzeReqInfo) (string, error) {
	if strings.Index(res.Url, "local:") == 0 {
		return "", nil
	}

	u, err := url.Parse(strings.ToLower(res.Url))
	if err != nil {
		return "", err
	}

	arr := strings.Split(u.Path, ".")
	if len(arr) > 1 {
		exname := arr[len(arr)-1]

		if exname == "mp3" || exname == "png" || exname == "gif" ||
			exname == "webp" || exname == "js" || exname == "css" ||
			exname == "mp4" {

			return exname, nil
		} else if exname == "jpg" || exname == "jpeg" {
			return "jpg", nil
		}
	}

	cct := strings.Split(res.ContentType, ";")
	cft := strings.Split(cct[0], "/")
	if len(cft) > 1 {
		if cft[1] == "jpg" || cft[1] == "jpeg" {
			cft[1] = "jpg"
		} else if cft[1] == "javascript" {
			cft[1] = "js"
		}

		return cft[1], nil
	}

	return cft[0], nil
}
