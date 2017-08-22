package youtube

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/otium/ytdl"
)

const (
	tmpPath = "youtube/tmp/"
)

func init() {
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}
}

type stream struct {
	token string
	title string
}

func (s *stream) DownloadVideo(url, customOutput string) {
	videoInfo, _ := ytdl.GetVideoInfo(url)
	s.title = videoInfo.Title

	videoInfo.Title = s.removeSpecialCharacter(videoInfo.Title)

	path := customOutput + videoInfo.Title + ".mp4"

	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	videoInfo.Download(videoInfo.Formats.Best(ytdl.FormatAudioEncodingKey)[0], file)
}

func (s *stream) removeSpecialCharacter(title string) string {
	chars := []string{"]", "^", "\\", "/", "'", "[", ".", "(", ")", "-", "?", "\""}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	title = re.ReplaceAllString(title, "")

	return title
}

func (s *stream) DownloadMP3(url, customOutput string) {
	s.token = s.randToken()
	tmp := tmpPath + s.token
	s.DownloadVideo(url, tmp)

	err := s.parseVideoToMP3(customOutput)
	if err != nil {
		fmt.Println(err)
	} else {
		s.removeTmpFile()
	}
}

func (s *stream) randToken() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *stream) parseVideoToMP3(customOutput string) error {
	path := s.getFilePath()
	fileName := s.getFile(path)
	output := customOutput + fileName

	_, err := exec.Command("ffmpeg", "-i", path, "-q:a", "0", "-map", "a", output).Output()

	return err
}

func (s *stream) getFilePath() string {
	path := ""
	files, _ := ioutil.ReadDir(tmpPath)
	for _, f := range files {
		if strings.Contains(f.Name(), s.token) {
			path = tmpPath + f.Name()
			break
		}
	}

	return path
}

func (s *stream) getFile(path string) string {
	file := ""
	values := strings.Split(path, s.token)
	file = values[1]

	return strings.Replace(file, "mp4", "mp3", -1)
}

func (s *stream) removeTmpFile() {
	path := s.getFilePath()
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}
