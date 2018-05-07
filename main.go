package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"image"
	"image/color"
	"image/gif"

	"github.com/caitlin615/gif-battle/kmeans"
)

var giphyAPIKey = os.Getenv("GIPHY_API_KEY")

func main() {
	if giphyAPIKey == "" {
		log.Fatal("missing Giphy API Key (use env var: `GIPHY_API_KEY`)")
	}

	// giphyGIF, err := randomGIPHYGIF("burrito")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// gifURL := giphyGIF.Images.DownsizedStill.URL
	gifURL := "https://media2.giphy.com/media/jtKpyFWJrVI08/giphy.gif"

	gifObject, err := GetGIF(gifURL)
	if err != nil {
		log.Fatal(err)
	}
	frame := gifObject.Image[0]

	mostUsedColor := getMostUsedColor(frame)
	log.Println("Most used color:", colorToHex(mostUsedColor))

	// using k-means clustering
	colors := getPixels(frame)
	km := kmeans.NewKmeansClustering(3, 5, 2)
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
	bounds := frame.Bounds()
	minX := bounds.Min.X
	minY := bounds.Min.Y
	maxX := bounds.Max.X
	maxY := bounds.Max.Y

	colors := map[color.Color]int{}

	// TODO: It's excessive to account for every pixel
	for xIndex := minX; xIndex < maxX; xIndex++ {
		for yIndex := minY; yIndex < maxY; yIndex++ {
			clr := frame.At(xIndex, yIndex)
			colors[clr]++
		}
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

// GIFRandom is the response from the random endpoint
type GIFRandom struct {
	Data GIPHYGIF `json:"data"`
	Meta Meta     `json:"meta"`
}

// Meta Object contains basic information regarding the request
type Meta struct {
	Message    string `json:"msg"`
	Status     int    `json:"status"`
	ResponseID string `json:"response_id"` // A unique ID paired with this response from the API.
}

// GIPHYGIF is a GIPHY Gif object: https://developers.giphy.com/docs/#gif-object
// Note: this doesn't include all values, only ones that seem necessary.
type GIPHYGIF struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Images Images `json:"images"`
}

// Images is the GIPHY Images object: https://developers.giphy.com/docs/#images-object
type Images struct {
	DownsizedStill Image `json:"downsized_still"` // will be used for finding prominent color
	Original       Image `json:"original"`
}

// Image is the GIPHY Image object
type Image struct {
	URL    string `json:"url"`
	Frames string `json:"frames"`
}

// GIF converts an Image into a GIF from the std lib image/gif package
func GetGIF(u string) (g *gif.GIF, err error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "image/gif")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return gif.DecodeAll(resp.Body)
}

func randomGIPHYGIF(keyword string) (gif GIPHYGIF, err error) {
	query := url.Values{}
	query.Add("api_key", giphyAPIKey)
	query.Add("tag", keyword)
	// query.Add("rating", "g")
	query.Add("fmt", "json")

	resp, err := http.Get("https://api.giphy.com/v1/gifs/random?" + query.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var giphyResp GIFRandom
	if err = json.Unmarshal(body, &giphyResp); err != nil {
		return
	}
	gif = giphyResp.Data

	// TODO: Check meta as well for status
	return
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
