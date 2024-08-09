package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/mcuadros/go-zcam-e2"
	"github.com/mcuadros/go-zcam-e2/settings"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cli := zcam.NewCamera(os.Getenv("CAMERA_IP"))

	if err := cli.StartSession(ctx); err != nil {
		log.Fatalf("error starting session: %s", err)
	}

	defer func() {
		if err := cli.QuitSession(ctx); err != nil {
			log.Fatalf("error quitting session: %s", err)
		}
	}()

	s := map[settings.Setting]any{
		settings.FlickerSetting:      "60Hz",
		settings.FocusSetting:        "AF",
		settings.ResolutionSetting:   "1920x1080",
		settings.VideoEncoderSetting: "H.265",
		settings.MovVFRSetting:       120,
	}

	log.Printf("configuring %d setting(s)", len(s))
	if err := cli.SetSettings(ctx, s); err != nil {
		log.Fatalf("error configuring camera: %s", err)
	}

	for setting := range s {
		log.Printf("reading setting %q", setting)
		value, err := cli.GetSetting(ctx, setting)
		if err != nil {
			log.Fatalf("error getting setting value %s: %s", setting, err)
		}

		fmt.Println(value)
	}

}
