package giphy

import (
	"encoding/json"
	"image/gif"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var giphyAPIKey = os.Getenv("GIPHY_API_KEY")

func init() {
	if giphyAPIKey == "" {
		log.Fatal("missing Giphy API Key (use env var: `GIPHY_API_KEY`)")
	}
}

// GIFRandom is the response from the random endpoint
type GIFRandom struct {
	Data GIF  `json:"data"`
	Meta Meta `json:"meta"`
}

// Meta Object contains basic information regarding the request
type Meta struct {
	Message    string `json:"msg"`
	Status     int    `json:"status"`
	ResponseID string `json:"response_id"` // A unique ID paired with this response from the API.
}

// GIF is a GIPHY Gif object: https://developers.giphy.com/docs/#gif-object
// Note: this doesn't include all values, only ones that seem necessary.
type GIF struct {
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

// GetGIF converts an Image into a GIF from the std lib image/gif package
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

// Random returns a GIF from GIPHY's random endpoint
func Random(keyword string) (gif GIF, err error) {
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
