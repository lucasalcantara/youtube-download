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

var wg sync.WaitGroup

func DownloadMP3(url string) {
	stream := stream{}
	stream.DownloadMP3(url)
}

func DownloadVideo(url string) {
	stream := stream{}
	stream.DownloadVideo(url, "")
}

func DownloadPlaylistVideos(url string) {
	downloadPlaylist(url, DownloadVideo)
}

func DownloadPlaylistMusics(url string) {
	downloadPlaylist(url, DownloadMP3)
}

func downloadPlaylist(url string, downloadFunc func(url string)) {
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

func startConcurrencyDownloadPlaylist(ids []string, downloadFunc func(url string)) {
	wg.Add(len(ids))

	for _, id := range ids {
		go func(videoId string, downloadFunc func(url string)) {
			log.Println("Start download for videoId: ", videoId)
			downloadFunc(youtubeUrl + videoId)
			log.Println("Finish download for videoId: ", videoId)

			wg.Done()
		}(id, downloadFunc)
	}

	wg.Wait()
}
