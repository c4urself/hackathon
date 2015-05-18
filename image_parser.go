package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	//"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
)

import _ "image/png"

type Region struct {
	img       image.Image
	offset    image.Point
	baseColor color.RGBA
}

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

func main() {
	fBaseImg, _ := os.Open("image.png")
	defer fBaseImg.Close()
	baseImg, _, err := image.Decode(fBaseImg)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var regions []Region = getImageRegions(baseImg, 40)
	log.Printf("Created %v regions", len(regions))

	for _, region := range regions {
		// useful debugging, remove later
		c := findBaseColor(region.img)
		region.baseColor = c

		fResultImg, _ := os.Create(fmt.Sprintf("%v.%v-%v-test.png", RGBToHex(c.R, c.G, c.B), region.offset.X, region.offset.Y))
		defer fResultImg.Close()
		png.Encode(fResultImg, region.img)
	}
	log.Printf("Generated base colours for each region")
}

func RGBToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func HexToRGB(h string) (uint8, uint8, uint8) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	if len(h) == 3 {
		h = h[:1] + h[:1] + h[1:2] + h[1:2] + h[2:] + h[2:]
	}
	if len(h) == 6 {
		if rgb, err := strconv.ParseUint(string(h), 16, 32); err == nil {
			return uint8(rgb >> 16), uint8((rgb >> 8) & 0xFF), uint8(rgb & 0xFF)
		}
	}
	return 0, 0, 0
}

func findBaseColor(img image.Image) color.RGBA {
	var rect image.Rectangle = img.Bounds()
	var length int = rect.Dx()

	var colorMap = make(map[string]int)

	// for each pixel add to a map
	for x := 0; x < length; x++ {
		for y := 0; y < length; y++ {
			var c color.RGBA = img.At(x, y).(color.RGBA)
			var hex string = RGBToHex(c.R, c.G, c.B)
			colorMap[hex] += 1
		}
	}

	// for each key, val in map return highest val
	var baseColor string
	var highest = 0
	for k, v := range colorMap {
		if v > highest {
			highest = v
			baseColor = k
		}
	}

	r, g, b := HexToRGB(baseColor)
	//log.Println("Found base color", baseColor)
	return color.RGBA{r, g, b, 255}
}
