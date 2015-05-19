package mosaic

import (
	"github.com/c4urself/hackathon/feeders"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
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

func MakeInstagramMosaic(username string, photosDir string, audienceDir string, mosaicDir string) {
	os.MkdirAll(photosDir, 0777)
	os.MkdirAll(audienceDir, 0777)
	os.MkdirAll(mosaicDir, 0777)

	feed := feeders.GetCreatorFeed(username)

	topPhotos := feed.Photos[:5]
	feeders.LoadPhotos(topPhotos, photosDir)
	feeders.LoadPhotos(feed.Audience, audienceDir)

	for _, photo := range topPhotos {
		photoPath := filepath.Join(photosDir, fmt.Sprintf("%s.png", photo.Id))
		mosaicPath := filepath.Join(mosaicDir, fmt.Sprintf("%s.png", photo.Id))
		MakeMosaic(photoPath, audienceDir, mosaicPath, 10)
	}
}
