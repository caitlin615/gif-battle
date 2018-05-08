package main

import (
	"encoding/json"
	"html/template"
	"image/color"
	"log"
	"net/http"
	"os"

	"github.com/caitlin615/gif-battle/colors"
	"github.com/caitlin615/gif-battle/giphy"
	"github.com/caitlin615/gif-battle/kmeans"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

type templateData struct {
	Left  battleImage
	Right battleImage
}

type battleImage struct {
	Image         *giphy.Image
	Colors        []color.Color
	MostUsedColor color.Color
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// "/" catches all paths, but we only want to execute the template on the root
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/new-battle", func(w http.ResponseWriter, r *http.Request) {
		data := templateData{
			Left:  newBattleGif(),
			Right: newBattleGif(),
		}

		err := json.NewEncoder(w).Encode(&data)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	port := "8080"
	if p, ok := os.LookupEnv("PORT"); ok && port != "" {
		port = p
	}
	addr := "0.0.0.0:" + port
	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func newBattleGif() (b battleImage) {
	giphyGIF, err := giphy.Random("burrito")
	if err != nil {
		log.Println(err)
		return
	}
	b.Image = &giphyGIF.Images.Original
	gifObject, err := giphy.GetGIF(giphyGIF.Images.DownsizedStill.URL)
	if err != nil {
		log.Println(err)
		return
	}
	frame := gifObject.Image[0]

	clrs := colors.GetPixels(frame)

	mostUsedColor := colors.GetMostUsedColor(frame)
	b.MostUsedColor = mostUsedColor

	km := kmeans.NewKmeansClustering(2, 3, 10)
	prominantColors := km.Run(clrs)

	for _, color := range prominantColors {
		b.Colors = append(b.Colors, color)
	}
	return
}
