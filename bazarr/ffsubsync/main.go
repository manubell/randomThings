package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var client = resty.New()

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	movies := getMovies()

	for _, movie := range movies.Data {
		var subtitlePath = ""
		for _, subtitle := range movie.Subtitles {
			if subtitle.Code3 == "nld" && len(subtitle.Path) > 0 {
				subtitlePath = subtitle.Path
			}
		}

		if subtitlePath != "" {
			syncDutchSubtitle(subtitlePath, movie.Path, movie.Title)
		}
	}
}

func syncDutchSubtitle(subtitlePath string, videoPath string, movieTitle string) {
	fmt.Println("Movie		:", movieTitle)
	resp, err := client.R().
		SetFormData(map[string]string{
			"apikey":        os.Getenv("BAZARR_API_KEY"),
			"language":      "nld",
			"subtitlesPath": subtitlePath,
			"videoPath":     videoPath,
			"mediaType":     "movies",
		}).Post(os.Getenv("BAZARR_HOST") + "/api/sync_subtitles")
	fmt.Println("Error      :", err)
	fmt.Println("Status     :", resp.Status())
	fmt.Println("Time       :", resp.Time())
	fmt.Println()
}

func getMovies() MovieList {
	response, err := client.R().
		SetQueryParams(map[string]string{
			"apikey":           os.Getenv("BAZARR_API_KEY"),
			"columns[0][data]": "monitored",
		}).Get(os.Getenv("BAZARR_HOST") + "/api/movies")

	if err != nil {
		log.Fatal(err)
	}

	movies := MovieList{}
	_ = json.Unmarshal(response.Body(), &movies)

	return movies
}

type MovieList struct {
	Data []Movie `json:"data"`
}

type Movie struct {
	Title     string      `json:"title"`
	Path      string      `json:"path"`
	Subtitles []Subtitles `json:"subtitles"`
}

type Subtitles struct {
	Code3 string `json:"code3"`
	Path  string `json:"path"`
}
