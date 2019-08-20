package main

import (
	"fmt"
	"io/ioutil"

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
}
