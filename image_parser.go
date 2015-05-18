package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func main() {

	type Region struct {
		img        image.Image
		baseColour string
		x          uint
		y          uint
	}

	fBaseImg, _ := os.Open("image.jpeg")
	defer fBaseImg.Close()
	baseImg, _, err := image.Decode(fBaseImg)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// image.RGBA is an in-memory image
	// image.Rect returns type Rectangle
	m := image.NewRGBA(image.Rect(0, 0, 20, 20))
	var rect image.Rectangle
	rect = m.Bounds()
	draw.Draw(m, rect, baseImg, image.Point{0, 0}, draw.Src)

	fResultImg, _ := os.Create("new.jpeg")
	defer fResultImg.Close()

	region := Region{img: m}
	jpeg.Encode(fResultImg, region.img, &jpeg.Options{jpeg.DefaultQuality})
}
