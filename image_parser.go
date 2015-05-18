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

// Split image on the set of regions
func getImageRegions(img image.Image, regionLength int) []Region {
	imgSize := img.Bounds().Size()
	xPointsCount := int(math.Ceil(float64(imgSize.X) / float64(regionLength)))
	yPointsCount := int(math.Ceil(float64(imgSize.Y) / float64(regionLength)))
	rgbImage := img.(*image.RGBA)

	var regions []Region;

	for i := 0; i < yPointsCount; i++ {
		for j := 0; j < xPointsCount; j++ {
			regionMinOffset := image.Point{X: j * regionLength, Y: i * regionLength}
			regionMaxOffset := image.Point{X: (j + 1) * regionLength, Y: (i + 1) * regionLength}
			regionImg := rgbImage.SubImage(image.Rectangle{Min: regionMaxOffset, Max: regionMaxOffset})
			regions = append(regions, Region{img: regionImg, offset: regionMinOffset})
		}
	}

	return regions
}

// Generate final image from the set of regions
func generateImage(path string, regions []Region) {
	// Get size of the final image
	width := 0;
	height := 0;

	for _, region := range regions {
		xOffset := region.offset.X + region.img.Bounds().Size().X
		yOffset := region.offset.Y + region.img.Bounds().Size().Y

		if width < xOffset {
			width = xOffset
		}

		if height < yOffset {
			height = yOffset
		}
	}

	// Draw the image
	destinationImage := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, region := range regions {
		draw.Draw(destinationImage,
				 image.Rectangle{Min: region.offset, Max: region.offset.Add(region.img.Bounds().Size())},
				 region.img,
				 region.offset,
				 draw.Src)
	}

	// Write it to file
	destinationFile, _ := os.Create(path)
	defer destinationFile.Close()
	png.Encode(destinationFile, destinationImage)
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
