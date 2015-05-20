package mosaic

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"math"
)

func GetColorProfile(img image.Image) color.NRGBA {
	// create a thumb to shortcut getting most common colours
	var length uint = 10
	var thumb image.Image = resize.Thumbnail(length, length, img, resize.NearestNeighbor)

	var r, g, b, count int64 = 0, 0, 0, 0

	// for each pixel add to a map
	for x := 0; x < int(length); x++ {
		for y := 0; y < int(length); y++ {
			c := color.NRGBAModel.Convert(thumb.At(x, y)).(color.NRGBA)
			r += int64(c.R)
			g += int64(c.G)
			b += int64(c.B)
			count += 1
		}
	}

	return color.NRGBA{uint8(r / count), uint8(g / count), uint8(b / count), 255}
}

func GetColorDistance(a color.NRGBA, b color.NRGBA) float64 {
	distance := math.Pow(float64(a.R)-float64(b.R), 2) +
		math.Pow(float64(a.G)-float64(b.G), 2) +
		math.Pow(float64(a.B)-float64(b.B), 2)

	return distance / (255 * 255 * 3)
}
