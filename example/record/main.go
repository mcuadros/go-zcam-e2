package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mcuadros/go-zcam-e2"
)

func main() {
	CameraIP := os.Getenv("CAMERA_IP")
	cli := zcam.NewCameraClient(fmt.Sprintf("http://%s", CameraIP))

	if err := cli.StartSession(); err != nil {
		log.Fatalf("error starting session: %s", err)
	}

	defer func() {
		if err := cli.QuitSession(); err != nil {
			log.Fatalf("error quitting session: %s", err)
		}
	}()

	log.Printf("starting recording for 1sec")
	if err := cli.StartVideoRecord(); err != nil {
		log.Fatalf("error starting record: %s", err)
	}

	time.Sleep(time.Second)

	log.Printf("stopping recording")
	if err := cli.StopVideoRecord(); err != nil {
		log.Fatalf("error starting record: %s", err)
	}

	files, err := cli.ListAllFiles()
	if err != nil {
		log.Fatalf("error listing files: %s", err)
	}

	for _, f := range files {
		log.Printf("downloading file: %s", f.Filename())
		bytes, err := f.Download(zcam.Original, f.Filename())
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("file %s downloaded, size %d bytes", f.Filename(), bytes)
		if err := f.Delete(); err != nil {
			log.Fatal(err)
		}
	}
}
