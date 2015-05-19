package mosaic

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	//"image/jpeg"
	"github.com/nfnt/resize"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const TILE_LENGTH = 10

type Region struct {
	img       image.Image
	offset    image.Point
	baseColor color.NRGBA
}

// Split image on the set of regions
func getImageRegions(img image.Image, regionLength int) []Region {
	imgSize := img.Bounds().Size()
	xPointsCount := int(math.Ceil(float64(imgSize.X) / float64(regionLength)))
	yPointsCount := int(math.Ceil(float64(imgSize.Y) / float64(regionLength)))
	rgbImage := img.(*image.RGBA)
	startPoint := image.Point{X: 0, Y: 0}

	log.Printf("Creating regions for %vx%v\n", xPointsCount, yPointsCount)

	var regions []Region

	for i := 0; i < yPointsCount; i++ {
		for j := 0; j < xPointsCount; j++ {
			// get the x,y start and x,y end for the region and make a rectangle
			regionMinOffset := image.Point{X: j * regionLength, Y: i * regionLength}
			regionMaxOffset := image.Point{X: (j + 1) * regionLength, Y: (i + 1) * regionLength}
			originalRect := image.Rectangle{Min: regionMinOffset, Max: regionMaxOffset}

			// create a rectangle based off of the old one with start at 0,0
			regionRect := image.Rectangle{startPoint, startPoint.Add(originalRect.Size())}
			regionImg := image.NewRGBA(regionRect)
			draw.Draw(regionImg, regionRect, rgbImage, regionMinOffset, draw.Src)

			// add to array of Regions
			region := Region{img: regionImg, offset: regionMinOffset}
			regions = append(regions, region)
		}
	}

	return regions
}

// Generate final image from the set of regions
func generateImage(path string, regions []Region) {
	// Get size of the final image
	width := 0
	height := 0

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
			image.Point{X: 0, Y: 0},
			draw.Src)
	}

	// Write it to file
	destinationFile, _ := os.Create(path)
	defer destinationFile.Close()
	png.Encode(destinationFile, destinationImage)
}

func theBestCandidate(region Region, candidates []Region) Region {
	var theBest Region
	var theBestSimilarity float64 = 0

	for _, candidate := range candidates {
		similarity := GetSimilarityOfColors(region.baseColor, candidate.baseColor)
		if similarity > theBestSimilarity {
			theBest.img = candidate.img
			theBest.offset = region.offset
			theBest.baseColor = candidate.baseColor
			theBestSimilarity = similarity
		}
	}
	return theBest
}

func matchRegions(originals []Region, thumbnails []image.Image) []Region {
	var populated []Region
	for _, region := range originals {
		populated = append(populated, Region{img: region.img, offset: region.offset, baseColor: FindBaseColor(region.img)})
	}

	var candidates []Region
	for _, img := range thumbnails {
		candidates = append(candidates, Region{img: img, baseColor: FindBaseColor(img)})
	}

	var result []Region
	for _, region := range populated {
		result = append(result, theBestCandidate(region, candidates))
	}

	return result
}

func generateImageSet(baseDir string) []image.Image {
	var imageSet []image.Image
	var totalFound int
	filepath.Walk(baseDir, func(path string, _ os.FileInfo, _ error) error {

		// only accept PNG for now
		if strings.HasSuffix(path, ".png") {
			totalFound += 1
			fImg, _ := os.Open(path)
			defer fImg.Close()
			img, _, err := image.Decode(fImg)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			m := resize.Resize(TILE_LENGTH, 0, img, resize.Lanczos3)
			imageSet = append(imageSet, m)
		}
		return nil
	})

	if totalFound == 0 {
		log.Panicf("No thumbnails found!")
	}

	log.Printf("Found %v thumbnails", totalFound)
	return imageSet
}
