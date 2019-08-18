package brucecore

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// GZipFile - gzip file infomation
type GZipFile struct {
	URL      string
	Size     int
	GZipSize int
}

func countGZipSize(in []byte) (int, error) {
	var buffer bytes.Buffer

	writer := gzip.NewWriter(&buffer)

	_, err := writer.Write(in)
	if err != nil {
		writer.Close()

		return 0, err
	}

	err = writer.Close()
	if err != nil {
		return 0, err
	}

	return len(buffer.Bytes()), nil
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	return data, err
}

// GetNoGZip - get no gzip resource
func GetNoGZip(reqs []*jarviscrawlercore.AnalyzeReqInfo) (
	[]GZipFile, error) {

	var lst []GZipFile

	for _, v := range reqs {
		rt, err := GetResType(v)
		if err != nil {
			return nil, err
		}

		if rt == "mp3" || rt == "webp" || rt == "gif" ||
			rt == "jpg" || rt == "png" || rt == "mp4" {

			continue
		}

		u, err := url.Parse(v.Url)
		if err != nil {
			continue
		}

		if u == nil {
			continue
		}

		lstpath := strings.Split(u.Path, "/")
		lstexn := strings.Split(lstpath[len(lstpath)-1], ".")
		if len(lstexn) <= 1 {
			continue
		}

		if v.ImgWidth == 0 && v.ImgHeight == 0 {
			if !v.IsGZip {
				buf, err := downloadFile(v.Url)
				if err == nil && buf != nil {
					size := len(buf)

					// exclude files smaller than 1k
					if size < 1024 {
						continue
					}

					nsize, err := countGZipSize(buf)
					if err == nil && nsize > 0 {
						lst = append(lst, GZipFile{
							URL:      v.Url,
							Size:     size,
							GZipSize: nsize,
						})
					}
				}
			}
		}
	}

	return lst, nil
}

// BuildNoGZipTableData - build [][]string
func BuildNoGZipTableData(files []GZipFile) [][]string {
	var lst [][]string
	ts := 0
	nts := 0

	for _, v := range files {
		lst = append(lst, []string{
			v.URL,
			FormatByteSize(v.Size),
			FormatByteSize(v.GZipSize),
			fmt.Sprintf("%.1f%%", float32(v.Size-v.GZipSize)*100/float32(v.Size)),
		})

		ts += v.Size
		nts += v.GZipSize
	}

	lst = append(lst, []string{
		"汇总",
		FormatByteSize(ts),
		FormatByteSize(nts),
		fmt.Sprintf("%.1f%%", float32(ts-nts)*100/float32(ts)),
	})

	return lst
}
