package brucecore

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/webp"

	"gonum.org/v1/gonum/stat"

	"github.com/lucasb-eyer/go-colorful"
	brucecorebase "github.com/zhs007/brucecore/base"
)

// DistanceColor - count distance c1 and c2
func DistanceColor(c1 color.Color, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	cf1 := colorful.Color{float64(r1>>8) / 255.0, float64(g1>>8) / 255.0, float64(b1>>8) / 255.0}
	cf2 := colorful.Color{float64(r2>>8) / 255.0, float64(g2>>8) / 255.0, float64(b2>>8) / 255.0}

	return cf1.DistanceCIE76(cf2)
}

// DistanceImg - distance image img1 and img2
func DistanceImg(fn1 string, fn2 string, mfn string) (image.Image, float64, error) {
	f1, err := os.Open(fn1)
	defer f1.Close()
	if err != nil {
		return nil, 0, err
	}

	img1, _, err := image.Decode(f1)
	if err != nil {
		return nil, 0, err
	}

	f2, err := os.Open(fn2)
	defer f2.Close()
	if err != nil {
		return nil, 0, err
	}

	img2, _, err := image.Decode(f2)
	if err != nil {
		return nil, 0, err
	}

	mf, err := os.Open(mfn)
	defer mf.Close()
	if err != nil {
		return nil, 0, err
	}

	mimg, _, err := image.Decode(mf)
	if err != nil {
		return nil, 0, err
	}

	if img1.Bounds().Size().X != img2.Bounds().Size().X || img1.Bounds().Size().Y != img2.Bounds().Size().Y {
		return nil, 0, brucecorebase.ErrInvalidImageRect
	}

	if img1.Bounds().Size().X != mimg.Bounds().Size().X || img1.Bounds().Size().Y != mimg.Bounds().Size().Y {
		return nil, 0, brucecorebase.ErrInvalidImageRect
	}

	// img := image.NewRGBA(image.Rect(0, 0, img1.Bounds().Size().X, img1.Bounds().Size().Y))
	img := image.NewRGBA(image.Rect(0, 0, img1.Bounds().Size().X, img1.Bounds().Size().Y))

	maxd := DistanceColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}, color.RGBA{R: 0, G: 0, B: 0, A: 0})

	var lst []float64

	for y := 0; y < img1.Bounds().Size().Y; y++ {
		for x := 0; x < img1.Bounds().Size().X; x++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)
			mc := mimg.At(x, y)

			d := DistanceColor(c1, c2)
			if d < 0 {
				d = -d
			}

			r1, g1, b1, _ := c1.RGBA()

			_, _, _, ma := mc.RGBA()

			// cf := colorful.Color{float64(r1>>8) / 255.0, float64(g1>>8) / 255.0, float64(b1>>8) / 255.0}
			// _, _, v1 := cf.Hsv()

			lst = append(lst, d*float64(ma>>8)/255.0)

			// r1, g1, b1, _ := c1.RGBA()

			img.Set(x, y, color.RGBA{
				R: uint8(r1 >> 8),
				G: uint8(g1 >> 8),
				B: uint8(b1 >> 8),
				A: uint8(d * 255.0 / maxd),
			})

			// img.Set(x, y, color.RGBA{
			// 	R: uint8(d * 255.0 / maxd),
			// 	G: uint8(d * 255.0 / maxd),
			// 	B: uint8(d * 255.0 / maxd),
			// 	A: uint8(d * 255.0 / maxd),
			// })
		}
	}

	v := stat.Variance(lst, nil)

	return img, v, nil
}

// GetSaturation - get saturation and variance
func GetSaturation(fn string) (image.Image, float64, error) {
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		return nil, 0, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, 0, err
	}

	// img := image.NewRGBA(image.Rect(0, 0, img1.Bounds().Size().X, img1.Bounds().Size().Y))
	imgret := image.NewGray(image.Rect(0, 0, img.Bounds().Size().X, img.Bounds().Size().Y))

	// maxd := DistanceColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}, color.RGBA{R: 0, G: 0, B: 0, A: 0})

	var lst []float64

	for y := 0; y < img.Bounds().Size().Y; y++ {
		for x := 0; x < img.Bounds().Size().X; x++ {
			c := img.At(x, y)

			r, g, b, _ := c.RGBA()

			cf := colorful.Color{float64(r>>8) / 255.0, float64(g>>8) / 255.0, float64(b>>8) / 255.0}
			_, s, _ := cf.Hsv()

			lst = append(lst, s)

			imgret.SetGray(x, y, color.Gray{
				Y: uint8(s * 255.0),
			})
		}
	}

	v := stat.Variance(lst, nil)

	return imgret, v, nil
}

// GetLightness - get lightness and variance
func GetLightness(fn string) (image.Image, float64, error) {
	f, err := os.Open(fn)
	defer f.Close()
	if err != nil {
		return nil, 0, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, 0, err
	}

	// img := image.NewRGBA(image.Rect(0, 0, img1.Bounds().Size().X, img1.Bounds().Size().Y))
	imgret := image.NewGray(image.Rect(0, 0, img.Bounds().Size().X, img.Bounds().Size().Y))

	// maxd := DistanceColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}, color.RGBA{R: 0, G: 0, B: 0, A: 0})

	var lst []float64

	for y := 0; y < img.Bounds().Size().Y; y++ {
		for x := 0; x < img.Bounds().Size().X; x++ {
			c := img.At(x, y)

			r, g, b, _ := c.RGBA()

			cf := colorful.Color{float64(r>>8) / 255.0, float64(g>>8) / 255.0, float64(b>>8) / 255.0}
			_, _, v := cf.Hsv()

			lst = append(lst, v)

			imgret.SetGray(x, y, color.Gray{
				Y: uint8(v * 255.0),
			})
		}
	}

	v := stat.Variance(lst, nil)

	return imgret, v, nil
}
