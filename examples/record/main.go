package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mcuadros/go-zcam-e2"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cli := zcam.NewCameraClient(os.Getenv("CAMERA_IP"))

	if err := cli.StartSession(ctx); err != nil {
		log.Fatalf("error starting session: %s", err)
	}

	defer func() {
		if err := cli.QuitSession(ctx); err != nil {
			log.Fatalf("error quitting session: %s", err)
		}
	}()

	log.Printf("starting recording for 1sec")
	f, err := cli.VideoRecord(ctx, time.Second)
	if err != nil {
		log.Fatalf("error recording: %s", err)
	}

	log.Printf("downloading file: %s", f.Filename())
	bytes, err := f.Download(ctx, zcam.Original, f.Filename())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("file %s downloaded, size %d bytes", f.Filename(), bytes)
	if err := f.Delete(ctx); err != nil {
		log.Fatal(err)
	}
}
