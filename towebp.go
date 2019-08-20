package brucecore

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/chai2010/webp"
)

// ToWebp - to webp
func ToWebp(fn string, Lossless bool, quality float32) ([]byte, error) {
	file, err := os.Open(fn)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = webp.Encode(&buf, img, &webp.Options{
		Lossless: Lossless,
		Quality:  quality,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
