package zcam

import (
	"context"
	"fmt"
	"strings"
)

type Stream string
type Setting string

const (
	// Stream0 by default, it's used as the file recording in stroage card.
	// The resolution and fps is controlled by 'movfmt' and 'movvfr'. The
	// encoder is set by 'video_encoder' in /ctrl/set interface.
	Stream0 Stream = "stream0"
	// Stream1 by default, it's used by the network streaming. The resolution
	// is limited by the stream 0, it can NOT be larger than the stream 0.
	// By default, the fps is 25 or 30, the encoder is H.264. The maxinum
	// resolution and fps is 4KP30.
	Stream1 Stream = "stream1"
)

const (
	//StreamSetting stream0/stream1
	StreamSetting Setting = "steam"
	// SettingWidth Video width in pixels",
	SettingWidth Setting = "width"
	// Video height in pixels",
	SettingHeight Setting = "height"
	// Encode bitrate in bps",
	SettingBitrate Setting = "bitrate"
	// Frames per second of the stream",
	SettingFPS Setting = "fps"
	// Video encoder (e.g., h264, h265)",
	SettingVenc Setting = "venc"
	// Bit width of the H.265 encoder",
	SettingBitwidth Setting = "bitwidth"
	// Designates the active stream for network streaming
	SettingSendStream Setting = "send_stream"
)

// SetStreamSource switches the stream source between internal options
func (c *Camera) SetStreamSource(ctx context.Context, stream Stream) error {
	endpoint := fmt.Sprintf("/ctrl/set?send_stream=%s", stream)
	return c.sendStreamRequest(ctx, endpoint)
}

// SetStreamSettings adjusts multiple settings for a designated stream
func (c *Camera) SetStreamSettings(ctx context.Context, stream Stream, settings map[Setting]string) error {
	if len(settings) == 0 {
		return fmt.Errorf("no settings provided")
	}

	// Construct query parameters from the settings map
	params := make([]string, 0, len(settings))
	for key, value := range settings {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}
	settingParams := strings.Join(params, "&")
	endpoint := fmt.Sprintf("/ctrl/stream_setting?index=%s&%s", stream, settingParams)

	return c.sendStreamRequest(ctx, endpoint)
}

type StreamConfig struct {
	Stream        Stream `json:"streamIndex"`
	EncoderType   string `json:"encoderType"`
	Bitwidth      string `json:"bitwidth"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	FPS           int    `json:"fps"`
	SampleUnit    int    `json:"sample_unit"`
	Bitrate       int    `json:"bitrate"`
	GopN          int    `json:"gop_n"`
	Rotation      int    `json:"rotation"`
	SplitDuration int    `json:"splitDuration"`
	Status        string `json:"status"`
}

// QueryStreamSetting retrieves the current settings for a specific stream
func (c *Camera) QueryStreamSetting(ctx context.Context, stream Stream) (*StreamConfig, error) {
	endpoint := fmt.Sprintf("/ctrl/stream_setting?index=%s&action=query", stream)
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var r StreamConfig
	if err := decodeJSON(body, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Helper function to send requests related to stream settings
func (c *Camera) sendStreamRequest(ctx context.Context, endpoint string) error {
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return decodeBasicRequest(body)
}
