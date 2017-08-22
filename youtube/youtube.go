package youtube

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"google.golang.org/api/youtube/v3"
)

const (
	listPrefix = "list="
	youtubeUrl = "https://www.youtube.com/watch?v="
)

func DownloadMP3(urls []string) {
	concurrencyDownload(urls, "", callStreamDownloadMP3)
}

func DownloadVideo(urls []string) {
	concurrencyDownload(urls, "", callStreamDownloadVideo)
}

func concurrencyDownload(urls []string, customOutput string, downloadFunc func(string, string)) {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url, customOutput string, downloadFunc func(string, string)) {
			downloadFunc(url, customOutput)
			wg.Done()
		}(url, customOutput, downloadFunc)
	}

	wg.Wait()
}

func callStreamDownloadVideo(url, customOutput string) {
	stream := stream{}
	stream.DownloadVideo(url, "")
}

func callStreamDownloadMP3(url, customOutput string) {
	stream := stream{}
	stream.DownloadMP3(url, "")
}

func DownloadPlaylistVideos(urls []string) {
	concurrentlyDownloadPlaylist(urls, DownloadVideo)
}

func DownloadPlaylistMusics(urls []string) {
	concurrentlyDownloadPlaylist(urls, DownloadMP3)
}

func concurrentlyDownloadPlaylist(urls []string, downloadFunc func(urls []string)) {
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			downloadPlaylist(url, downloadFunc)
			wg.Done()
		}(url)
	}

	wg.Wait()
}

func downloadPlaylist(url string, downloadFunc func(url []string)) {
	service := youtubeService()
	listId := playlistId(url)
	ids := playlistItemsIds(listId, service)

	startConcurrencyDownloadPlaylist(ids, downloadFunc)

	fmt.Println("Playlist Downloaded")
}

func youtubeService() *youtube.Service {
	client := getClient()
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

func startConcurrencyDownloadPlaylist(ids []string, downloadFunc func(url []string)) {
	var wg sync.WaitGroup
	wg.Add(len(ids))

	for _, id := range ids {
		go func(videoId string, downloadFunc func(url []string)) {
			log.Println("Start download for videoId: ", videoId)
			url := []string{youtubeUrl + videoId}

			downloadFunc(url)
			log.Println("Finish download for videoId: ", videoId)

			wg.Done()
		}(id, downloadFunc)
	}

	wg.Wait()
}
