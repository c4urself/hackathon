package mosaic

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
)

func GetGolorPorfile(img image.Image) color.NRGBA {
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

func GetColorDistance(a color.NRGBA, b color.NRGBA) float64 {
	distance := math.Pow(float64(a.R)-float64(b.R), 2) +
		math.Pow(float64(a.G)-float64(b.G), 2) +
		math.Pow(float64(a.B)-float64(b.B), 2)

	return distance / (255*255*3)
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
