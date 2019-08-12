package brucecore

import (
	"fmt"
	"strings"

	"net/url"
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

// GetSource - https://www.a.com/b/c.png => https://www.a.com
func GetSource(str string) (string, error) {
	if strings.Index(str, "local:imgdata-") == 0 {
		return "local:imgdata", nil
	}

	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil
}
