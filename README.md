# go-zcam-e2 [![GoDoc](https://godoc.org/github.com/mcuadros/go-zcam-e2?status.svg)](https://godoc.org/github.com/mcuadros/go-zcam-e2) [![GitHub release](https://img.shields.io/github/release/mcuadros/go-zcam-e2.svg)](https://github.com/mcuadros/go-zcam-e2/releases)
==============================

The `go-zcam-e2` provides a comprehensive interface for interacting with [Z CAM](https://www.z-cam.com/) E2 series. It supports various functionalities including focus and zoom control, file management, card management, and network streaming.

Features
--------

- Focus & Zoom Control: Manage autofocus, manual focus adjustments, and zoom functionalities directly through HTTP commands.
- File Management: List files, download, delete, and retrieve metadata for files stored on the camera.
- Card Management: Check card presence, format the storage card, and query storage information.
- Network Streaming: Manage streaming settings, switch between different streams, and configure streaming parameters like resolution and bitrate.

Prerequisites
-------------

- Go (Golang) environment setup (Version 1.14 or higher recommended)
- Network access to a Z CAM E2 camera


Installation
------------

The recommended way to install go-zcam-e2

```
go get github.com/mcuadros/go-zcam-e2
```

Examples
--------


```go
import "github.com/mcuadros/go-zcam-e2"
```

```go
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
```

License
-------

MIT, see [LICENSE](LICENSE)
