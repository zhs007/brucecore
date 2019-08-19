package brucecore

import (
	"fmt"

	"github.com/zhs007/adacore"
	jarviscrawlercore "github.com/zhs007/jccclient/proto"
)

// GetImageRectName - (100, 50) => 100x50, (50, 100) => 100x50
func GetImageRectName(w int, h int) string {
	w1 := w
	h1 := h
	if w1 < h1 {
		w1 = h
		h1 = w
	}

	return fmt.Sprintf("%vx%v", w1, h1)
}

// Image - image
type Image struct {
	URL          string
	ImgType      string
	Width        int
	Height       int
	Size         int
	PixelQuality float32
}

// ImageMgr - image manager
type ImageMgr struct {
	Imgs     []*Image
	TreeMap  []adacore.ChartTreeMapSeriesNodeFloat
	HostList []string
}

// Insert - insert a SourceInfo
func (mgr *ImageMgr) Insert(url string, imgtype string, bytes int, w int, h int) {
	v := mgr.Find(url)
	if v == nil {
		mgr.Imgs = append(mgr.Imgs, &Image{
			URL:          url,
			ImgType:      imgtype,
			Width:        w,
			Height:       h,
			Size:         bytes,
			PixelQuality: float32(bytes) / float32(w*h),
		})

		return
	}
}

// Find - find a SourceInfo
func (mgr *ImageMgr) Find(url string) *Image {
	for _, v := range mgr.Imgs {
		if v.URL == url {
			return v
		}
	}

	return nil
}

// buildTreeMapData - build treemap data
func (mgr *ImageMgr) buildTreeMapData() {
	mgr.HostList = []string{"images"}
	mgr.TreeMap = []adacore.ChartTreeMapSeriesNodeFloat{
		adacore.ChartTreeMapSeriesNodeFloat{
			Name: "images",
		},
	}

	for _, v := range mgr.Imgs {
		mgr.insertImageType2TreeMap(v)
	}

	for _, v := range mgr.Imgs {
		mgr.insertImageRect2TreeMap(v)
	}

	for _, v := range mgr.Imgs {
		mgr.insertImage2TreeMap(v)
	}
}

// insertImageType2TreeMap - insert image type
func (mgr *ImageMgr) insertImageType2TreeMap(img *Image) {
	for _, v := range mgr.TreeMap[0].Data {
		if v.Name == img.ImgType {
			return
		}
	}

	mgr.TreeMap[0].Data = append(mgr.TreeMap[0].Data, adacore.ChartTreeMapDataFloat{
		Name: img.ImgType,
	})
}

// insertImageRect2TreeMap - insert image rect
func (mgr *ImageMgr) insertImageRect2TreeMap(img *Image) {
	strname := GetImageRectName(img.Width, img.Height)

	for i := range mgr.TreeMap[0].Data {
		if mgr.TreeMap[0].Data[i].Name == img.ImgType {
			for _, rv := range mgr.TreeMap[0].Data[i].Children {
				if rv.Name == strname {
					return
				}
			}

			mgr.TreeMap[0].Data[i].Children = append(mgr.TreeMap[0].Data[i].Children, adacore.ChartTreeMapDataFloat{
				Name: strname,
			})

			return
		}
	}
}

// insertImage - insert image
func (mgr *ImageMgr) insertImage2TreeMap(img *Image) {
	strname := GetImageRectName(img.Width, img.Height)

	for i := range mgr.TreeMap[0].Data {
		if mgr.TreeMap[0].Data[i].Name == img.ImgType {
			for j := range mgr.TreeMap[0].Data[i].Children {
				if mgr.TreeMap[0].Data[i].Children[j].Name == strname {
					mgr.TreeMap[0].Data[i].Children[j].Children = append(
						mgr.TreeMap[0].Data[i].Children[j].Children,
						adacore.ChartTreeMapDataFloat{
							Name:  GetURLName(img.URL),
							Value: img.PixelQuality,
							URL:   img.URL,
						})

					return
				}
			}
		}
	}
}

// AnalyzeImageMgr - analyze image manager
func AnalyzeImageMgr(reqs []*jarviscrawlercore.AnalyzeReqInfo) (*ImageMgr, error) {
	mgr := &ImageMgr{}

	for _, v := range reqs {
		rt, err := GetResType(v)
		if err == nil && rt != "" && (v.ImgWidth >= 64 || v.ImgHeight >= 64) && v.BufBytes >= 1024 {

			mgr.Insert(v.Url, rt, int(v.BufBytes), int(v.ImgWidth), int(v.ImgHeight))
		}
	}

	mgr.buildTreeMapData()

	return mgr, nil
}
