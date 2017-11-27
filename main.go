package main

import (
	//"drive"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"youtube"
)

type youtubeFunc func(url []string, upload bool)

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
		urls, _ := url.ParseQuery(r.Form.Get("urls"))
		format := r.Form.Get("format")
		upload, err := strconv.ParseBool(r.Form.Get("upload"))
		if err != nil {
			log.Fatal(err)
		}

		option := strings.ToLower(r.Form.Get("option"))

		youtubeFuncs[option][format](urls["url"], upload)

		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {

	//drive.Print()

	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/post", postHandler)

	log.Println("Listening...")
	http.ListenAndServe(":80", nil)
}
