package colors

import (
	"fmt"
	"image"
	"image/color"
	"log"
)

// GetPixels returns all the colors included in the supplied image
func GetPixels(frame *image.Paletted) (colors []color.Color) {
	bounds := frame.Bounds()
	minX := bounds.Min.X
	minY := bounds.Min.Y
	maxX := bounds.Max.X
	maxY := bounds.Max.Y

	// TODO: It's excessive to account for every pixel
	for xIndex := minX; xIndex < maxX; xIndex++ {
		for yIndex := minY; yIndex < maxY; yIndex++ {
			clr := frame.At(xIndex, yIndex)
			colors = append(colors, clr)
		}
	}
	return
}

// GetMostUsedColor returns the color that occurs most often in the image
func GetMostUsedColor(frame *image.Paletted) color.Color {
	colors := map[color.Color]int{}
	for _, clr := range GetPixels(frame) {
		colors[clr]++
	}

	if len(colors) == 0 {
		log.Fatal("no colors")
	}

	// Find the color with the most number of occurances
	// TODO: Sort might be better
	maxCount := 0
	var maxColor color.Color
	for color, count := range colors {
		if count > maxCount {
			maxCount = count
			maxColor = color
		}
	}
	return maxColor
}

// ColorToHex converts a color.Color to a hex string
// https://github.com/pwaller/go-hexcolor
func ColorToHex(c color.Color) string {
	r, g, b, a := c.RGBA()
	return RGBAToHex(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

// RGBAToHex converts an RGBA to a Hex string.
// If a == 255, the A is not specified in the hex string
func RGBAToHex(r, g, b, a uint8) string {
	if a == 255 {
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a)
}
