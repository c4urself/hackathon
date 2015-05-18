package main

import (
	"fmt"
	"image"
	"os"
	"math"
)
import _ "image/png"

type Region struct {
	img image.Image
	offset image.Point
	baseColour string
}

func getImageRegions(img image.Image, regionLength int) []Region {
	var imgSize image.Point = img.Bounds().Size()
	var xPointsCount int = int(math.Ceil(float64(imgSize.X) / float64(regionLength)))
	var yPointsCount int = int(math.Ceil(float64(imgSize.Y) / float64(regionLength)))
	var regions []Region;
	var regionImg image.Image;
	var regionMinOffset image.Point;
	var regionMaxOffset image.Point;

	rgbImage := img.(*image.RGBA)

	for i := 0; i < yPointsCount; i++ {
		for j := 0; j < xPointsCount; j++ {
			regionMinOffset = image.Point{X: j * regionLength, Y: i * regionLength}
			regionMaxOffset = image.Point{X: (j + 1) * regionLength, Y: (i + 1) * regionLength}
			regionImg = rgbImage.SubImage(image.Rectangle{Min: regionMaxOffset, Max: regionMaxOffset})
			regions = append(regions, Region{img: regionImg, offset: regionMinOffset})
		}
	}

	return regions
}

func main() {
	fBaseImg, _ := os.Open("image.png")
	defer fBaseImg.Close()
	baseImg, _, err := image.Decode(fBaseImg)
	
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(getImageRegions(baseImg, 40));
}
