package main

import (
	"fmt"
	"log"

	"image"
	"image/color"

	"github.com/caitlin615/gif-battle/giphy"
	"github.com/caitlin615/gif-battle/kmeans"
)

func main() {
	// giphyGIF, err := giphy.Random("burrito")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// gifURL := giphyGIF.Images.DownsizedStill.URL
	gifURL := "https://media2.giphy.com/media/jtKpyFWJrVI08/giphy.gif"

	gifObject, err := giphy.GetGIF(gifURL)
	if err != nil {
		log.Fatal(err)
	}
	frame := gifObject.Image[0]

	// using most common color
	mostUsedColor := getMostUsedColor(frame)
	log.Println("Most used color:", colorToHex(mostUsedColor))

	// using k-means clustering
	colors := getPixels(frame)
	km := kmeans.NewKmeansClustering(2, 10, 10)
	prominantColors := km.Run(colors)
	for _, color := range prominantColors {
		log.Println(colorToHex(color))
	}
}

func getPixels(frame *image.Paletted) (colors []color.Color) {
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

// getMostUsedColor returns the color that occurs most often in the image
func getMostUsedColor(frame *image.Paletted) color.Color {
	colors := map[color.Color]int{}
	for _, clr := range getPixels(frame) {
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

// https://github.com/pwaller/go-hexcolor
func colorToHex(c color.Color) string {
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
