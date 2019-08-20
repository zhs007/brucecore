package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/zhs007/brucecore"
)

func img2webp(fn string, outputfn string, Lossless bool, q float32) error {
	buf, err := brucecore.ToWebp(fn, Lossless, q)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputfn, buf, 0644)
	if err != nil {
		return err
	}

	return nil
}

func distanceimg(fn1 string, fn2 string, ofn string) error {
	img, v, err := brucecore.DistanceImg(fn1, fn2)
	if err != nil {
		return err
	}

	fmt.Printf("v = %v", v)

	f, _ := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)

	return nil
}

func main() {
	err := img2webp("./test001.png", fmt.Sprintf("./test001-l.webp"), true, 0)
	if err != nil {
		fmt.Printf("img2webp error %v", err)
	}

	for i := 100; i > 0; i -= 10 {
		err := img2webp("./test001.png", fmt.Sprintf("./test001-q%v.webp", i), false, float32(i))
		if err != nil {
			fmt.Printf("img2webp error %v", err)
		}
	}

	err = distanceimg("./test001.png", "./test001-q10.webp", "./output.png")
	if err != nil {
		fmt.Printf("img2webp error %v", err)
	}
}
