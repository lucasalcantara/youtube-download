package youtube

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/otium/ytdl"
)

var customOutput = ""

const tmpPath = "youtube/tmp/"

func init() {
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}
}

type stream struct {
	title string
}

func (s *stream) DownloadVideo(url string) string {
	videoInfo, _ := ytdl.GetVideoInfo(url)

	videoInfo.Title = s.removeSpecialCharacter(videoInfo.Title)
	s.title = videoInfo.Title

	p := customOutput + videoInfo.Title + ".mp4"

	file, err := os.Create(p)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	videoInfo.Download(videoInfo.Formats.Best(ytdl.FormatAudioEncodingKey)[0], file)
	return p
}

func (s *stream) removeSpecialCharacter(title string) string {
	chars := []string{"]", "^", "\\", "/", "'", "[", ".", "(", ")", "-", "?", "\""}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	title = re.ReplaceAllString(title, "")

	return title
}

func (s *stream) DownloadMP3(url string) string {
	customOutput = tmpPath
	s.DownloadVideo(url)
	customOutput = ""

	p, err := s.parseVideoToMP3()
	if err != nil {
		log.Println(err)
	} else {
		s.removeTmpFile()
	}

	return p
}

func (s *stream) parseVideoToMP3() (string, error) {
	p := s.getFilePath()
	fileName := s.getFileName(p)

	_, err := exec.Command("ffmpeg", "-i", p, "-q:a", "0", "-map", "a", fileName).Output()

	return fileName, err
}

func (s *stream) getFilePath() string {
	p := ""

	files, _ := ioutil.ReadDir(tmpPath)
	for _, f := range files {
		if strings.Contains(f.Name(), s.title) {
			p = tmpPath + f.Name()
			break
		}
	}

	return p
}

func (s *stream) getFileName(p string) string {
	filename := path.Base(p)
	return strings.Replace(filename, "mp4", "mp3", -1)
}

func (s *stream) removeTmpFile() {
	p := s.getFilePath()
	err := os.Remove(p)
	if err != nil {
		log.Println(err)
	}
}
