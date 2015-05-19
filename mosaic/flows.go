package mosaic

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func MakeMosaic(mainPath string, thumbnailsPath string, mosaicPath string, tileSize int) {
	// Import base image
	baseImgFile, _ := os.Open(mainPath)
	defer baseImgFile.Close()
	baseImg, _, err := image.Decode(baseImgFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Generate mosaic
	originalTiles := BreakToTiles(baseImg, tileSize)
	candidateTiles := ImportTiles(thumbnailsPath, tileSize)
	similarTiles := FindSimilarTiles(originalTiles, candidateTiles) 
	mosaicImage := CollectFromTiles(similarTiles)

	// Write result
	mosaicFile, _ := os.Create(mosaicPath)
	defer mosaicFile.Close()
	png.Encode(mosaicFile, mosaicImage)
}
