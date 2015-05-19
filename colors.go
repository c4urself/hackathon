package mosaic

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
)

// Determines similarity of colors (from 0 to 1). Higher - better
func GetSimilarityOfColors(f color.NRGBA, s color.NRGBA) float64 {
	distance := math.Pow(float64(f.R)-float64(s.R), 2) +
		math.Pow(float64(f.G)-float64(s.G), 2) +
		math.Pow(float64(f.B)-float64(s.B), 2)

	return 1 - distance/(255*255*3)
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

func FindBaseColor(img image.Image) color.NRGBA {
	var rect image.Rectangle = img.Bounds()
	var length int = rect.Dx()

	var colorMap = make(map[string]int)

	// for each pixel add to a map
	for x := 0; x < length; x++ {
		for y := 0; y < length; y++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
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
	return color.NRGBA{r, g, b, 255}
}
