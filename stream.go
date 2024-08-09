package zcam

import (
	"context"
	"fmt"
	"strings"
)

type Stream string
type Setting string

const (
	Stream0 Stream = "stream0"
	Stream1 Stream = "stream1"
)

type StreamSettingResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}

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
func (c *Camera) SetStreamSource(ctx context.Context, stream Stream) (*StreamSettingResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/set?send_stream=%s", stream)
	return c.sendStreamRequest(ctx, endpoint)
}

// SetStreamSettings adjusts multiple settings for a designated stream
func (c *Camera) SetStreamSettings(ctx context.Context, stream Stream, settings map[Setting]string) (*StreamSettingResponse, error) {
	if len(settings) == 0 {
		return nil, fmt.Errorf("no settings provided")
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

// QueryStreamSetting retrieves the current settings for a specific stream
func (c *Camera) QueryStreamSetting(ctx context.Context, stream Stream) (*StreamSettingResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/stream_setting?index=%s&action=query", stream)
	return c.sendStreamRequest(ctx, endpoint)
}

// Helper function to send requests related to stream settings
func (c *Camera) sendStreamRequest(ctx context.Context, endpoint string) (*StreamSettingResponse, error) {
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var response StreamSettingResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
