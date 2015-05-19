package mosaic

import (
	"fmt"
	"image"
	"log"
	"os"
)

func MakeMosaic(mainPath string, thumbnailsPath string, mosaicPath string) {
	fBaseImg, _ := os.Open(mainPath)
	defer fBaseImg.Close()
	baseImg, _, err := image.Decode(fBaseImg)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var regions []Region = getImageRegions(baseImg, 40)
	log.Printf("Created %v regions", len(regions))

	result := matchRegions(regions, generateImageSet(thumbnailsPath))
	generateImage(mosaicPath, result)
}
