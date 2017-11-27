package drive

import (
	"common"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
)

const driveQuerySpace = "drive-credentials.json"

func UploadFile(path string) {
	f := file(path)

	defer removePath(path)
	defer f.Close()

	srv := server()

	_, err := srv.Files.Create(&drive.File{Name: f.Name()}).Media(f).Do()
	if err != nil {
		log.Fatalf("Unable to create file: %v Path: %v", err, path)
	}
}

func file(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening %q: %v", path, err)
	}

	return f
}

func server() *drive.Service {
	client := common.GetClient(drive.DriveFileScope, driveQuerySpace)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	return srv
}

func removePath(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("error removing %q: %v", path, err)
	}
}
