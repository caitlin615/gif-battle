package main

import (
	"log"

	"github.com/caitlin615/gif-battle/colors"
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
	mostUsedColor := colors.GetMostUsedColor(frame)
	log.Println("Most used color:", colors.ColorToHex(mostUsedColor))

	// using k-means clustering
	clrs := colors.GetPixels(frame)
	km := kmeans.NewKmeansClustering(2, 10, 10)
	prominantColors := km.Run(clrs)
	for _, color := range prominantColors {
		log.Println(colors.ColorToHex(color))
	}
}
