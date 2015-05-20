package mosaic

import (
	"fmt"
	"github.com/c4urself/hackathon/feeders"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

type Mosaic struct {
	Id          string
	OriginalUrl string
	RelativeUrl string
	Likes       int64
}

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

func MakeInstagramMosaic(username string, photosDir string, audienceDir string, mosaicDir string) []Mosaic {
	os.RemoveAll(photosDir)
	os.RemoveAll(audienceDir)
	os.RemoveAll(mosaicDir)

	os.MkdirAll(photosDir, 0777)
	os.MkdirAll(audienceDir, 0777)
	os.MkdirAll(mosaicDir, 0777)

	feed := feeders.GetCreatorFeed(username)
	topPhotos := feed.GetTopPhotos(6)

	feeders.LoadPhotos(topPhotos, photosDir)
	feeders.LoadPhotos(feed.Audience, audienceDir)

	var mosaic []Mosaic

	for _, photo := range topPhotos {
		photoPath := filepath.Join(photosDir, fmt.Sprintf("%s.png", photo.Id))
		mosaicPath := filepath.Join(mosaicDir, fmt.Sprintf("%s.png", photo.Id))
		MakeMosaic(photoPath, audienceDir, mosaicPath, 5)

		mosaic = append(mosaic, Mosaic{
			Id:          photo.Id,
			RelativeUrl: fmt.Sprintf("%s.png", photo.Id),
			OriginalUrl: photo.Url,
			Likes:       photo.Likes})
	}

	return mosaic
}
