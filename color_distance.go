package main

import (
	"fmt"
	"image/color"
	"math"
)

// Determines similarity of colors (from 0 to 1). Higher - better
func getOldDistanceBetweenColors(first color.Color, second color.Color) float64 {
	f := first.(color.RGBA)
	s := second.(color.RGBA)

	distance := math.Pow(float64(f.R)-float64(s.R), 2) +
		math.Pow(float64(f.G)-float64(s.G), 2) +
		math.Pow(float64(f.B)-float64(s.B), 2)

	return 1 - distance/(255*255*3)
}

func oldMain() {
	fmt.Println(getDistanceBetweenColors(color.RGBA{255, 51, 51, 1}, color.RGBA{255, 153, 153, 1}))
	fmt.Println(getDistanceBetweenColors(color.RGBA{255, 51, 51, 1}, color.RGBA{102, 178, 255, 1}))
}
