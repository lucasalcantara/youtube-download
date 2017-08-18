package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"youtube"
)

type youtubeFunc func(url string)

var (
	youtubeFuncs map[string]map[string]youtubeFunc
)

func init() {
	youtubeFuncs = make(map[string]map[string]youtubeFunc)

	youtubeFuncs["video"] = make(map[string]youtubeFunc)
	youtubeFuncs["playlist"] = make(map[string]youtubeFunc)

	youtubeFuncs["video"]["MP3"] = youtube.DownloadMP3
	youtubeFuncs["video"]["MP4"] = youtube.DownloadVideo

	youtubeFuncs["playlist"]["MP3"] = youtube.DownloadPlaylistMusics
	youtubeFuncs["playlist"]["MP4"] = youtube.DownloadPlaylistVideos

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		url := r.Form.Get("url")
		format := r.Form.Get("format")
		option := strings.ToLower(r.Form.Get("option"))

		youtubeFuncs[option][format](url)

		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/post", postHandler)

	log.Println("Listening...")
	http.ListenAndServe(":80", nil)
}
