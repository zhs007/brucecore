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

	cf1 := colorful.Color{float64(r1) / 255.0, float64(g1) / 255.0, float64(b1) / 255.0}
	cf2 := colorful.Color{float64(r2) / 255.0, float64(g2) / 255.0, float64(b2) / 255.0}

	return cf1.DistanceCIE76(cf2)
}

// DistanceImg - distance image img1 and img2
func DistanceImg(fn1 string, fn2 string) (image.Image, float64, error) {
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

	if img1.Bounds().Size().X != img2.Bounds().Size().X || img1.Bounds().Size().Y != img2.Bounds().Size().Y {
		return nil, 0, brucecorebase.ErrInvalidImageRect
	}

	// img := image.NewRGBA(image.Rect(0, 0, img1.Bounds().Size().X, img1.Bounds().Size().Y))
	img := image.NewGray(image.Rect(0, 0, img1.Bounds().Size().X, img1.Bounds().Size().Y))

	maxd := DistanceColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}, color.RGBA{R: 0, G: 0, B: 0, A: 0})

	var lst []float64

	for y := 0; y < img1.Bounds().Size().Y; y++ {
		for x := 0; x < img1.Bounds().Size().X; x++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)

			d := DistanceColor(c1, c2)
			if d < 0 {
				d = -d
			}

			lst = append(lst, d/maxd)

			img.SetGray(x, y, color.Gray{Y: uint8(d / maxd * 255.0)})
		}
	}

	v := stat.Variance(lst, nil)

	return img, v, nil
}
