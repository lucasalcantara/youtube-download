package youtube

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"

	"common"
	"drive"

	"google.golang.org/api/youtube/v3"
)

var (
	download chan params
	wg       sync.WaitGroup
)

const (
	listPrefix        = "list="
	youtubeUrl        = "https://www.youtube.com/watch?v="
	youtubeQuerySpace = "youtube-credentials.json"
)

type params struct {
	downloadFunc func(string) string
	url          string
	upload       bool
}

func init() {
	download = make(chan params, 1)

	for i := 0; i < runtime.NumCPU(); i++ {
		go run()
	}
}

func DownloadMP3(urls []string, upload bool) {
	concurrencyDownload(urls, upload, callStreamDownloadMP3)
}

func DownloadVideo(urls []string, upload bool) {
	concurrencyDownload(urls, upload, callStreamDownloadVideo)
}

func run() {
	for {
		select {
		case param := <-download:
			path := param.downloadFunc(param.url)
			if path != "" && param.upload {
				drive.UploadFile(path)
			}

			wg.Done()
		}
	}
}

func concurrencyDownload(urls []string, upload bool, downloadFunc func(string) string) {
	wg.Add(len(urls))

	for _, url := range urls {
		download <- params{downloadFunc, url, upload}
	}

	wg.Wait()
}

func callStreamDownloadVideo(url string) string {
	stream := stream{}
	return stream.DownloadVideo(url)
}

func callStreamDownloadMP3(url string) string {
	stream := stream{}
	return stream.DownloadMP3(url)
}

func DownloadPlaylistVideos(urls []string, upload bool) {
	concurrentlyDownloadPlaylist(urls, upload, DownloadVideo)
}

func DownloadPlaylistMusics(urls []string, upload bool) {
	concurrentlyDownloadPlaylist(urls, upload, DownloadMP3)
}

func concurrentlyDownloadPlaylist(urls []string, upload bool, downloadFunc func(urls []string, upload bool)) {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			downloadPlaylist(url, upload, downloadFunc)
			wg.Done()
		}(url)
	}

	wg.Wait()
}

func downloadPlaylist(url string, upload bool, downloadFunc func(urls []string, upload bool)) {
	service := youtubeService()
	listId := playlistId(url)
	ids := playlistItemsIds(listId, service)
	urls := make([]string, 0)

	for _, id := range ids {
		urls = append(urls, youtubeUrl+id)
	}

	downloadFunc(urls, upload)

	fmt.Println("Playlist Downloaded")
}

func youtubeService() *youtube.Service {
	client := common.GetClient(youtube.YoutubeReadonlyScope, youtubeQuerySpace)
	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error to get the service.")
	}

	return service
}

func playlistId(url string) string {
	id := ""
	getParam := strings.Split(url, "?")[1]
	params := strings.Split(getParam, "&")

	for _, param := range params {
		if strings.Contains(param, listPrefix) {
			id = strings.Replace(param, listPrefix, "", -1)
			break
		}
	}

	return id
}

func playlistItemsIds(playListId string, service *youtube.Service) []string {
	ids := make([]string, 0)
	nextPageToken := ""
	for {
		playlistResponse := playlistItemsList(service, "snippet", playListId, nextPageToken)

		for _, playlistItem := range playlistResponse.Items {
			ids = append(ids, playlistItem.Snippet.ResourceId.VideoId)
		}

		nextPageToken = playlistResponse.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return ids
}

func playlistItemsList(service *youtube.Service, part string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	call := service.PlaylistItems.List(part)
	call = call.PlaylistId(playlistId)

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error getting playlistitems: %v", err)
	}

	return response
}
