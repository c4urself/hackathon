package main

import (
	"fmt"
	"image/color"
	"math"
)

// Determines similarity of colors (from 0 to 1)
func getDistanceBetweenColors(first color.Color, second color.Color) float64 {
	f := first.(color.RGBA)
	s := second.(color.RGBA)

	return math.Pow(float64(f.R) - float64(s.R), 2) + 
	       math.Pow(float64(f.G) - float64(s.G), 2) +
	       math.Pow(float64(f.B) - float64(s.B), 2);
}
