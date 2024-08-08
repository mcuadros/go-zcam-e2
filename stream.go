package zcam

import (
	"fmt"
	"strings"
)

// SetStreamSource switches the stream source between internal options
func (c *CameraClient) SetStreamSource(stream string) (*StreamSettingResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/set?send_stream=%s", stream)
	return c.sendStreamRequest(endpoint)
}

// SetStreamSettings adjusts multiple settings for a designated stream
func (c *CameraClient) SetStreamSettings(index string, settings map[string]string) (*StreamSettingResponse, error) {
	if len(settings) == 0 {
		return nil, fmt.Errorf("no settings provided")
	}

	// Construct query parameters from the settings map
	params := make([]string, 0, len(settings))
	for key, value := range settings {
		params = append(params, fmt.Sprintf("%s=%s", key, value))
	}
	settingParams := strings.Join(params, "&")
	endpoint := fmt.Sprintf("/ctrl/stream_setting?index=%s&%s", index, settingParams)

	return c.sendStreamRequest(endpoint)
}

// QueryStreamSetting retrieves the current settings for a specific stream
func (c *CameraClient) QueryStreamSetting(index string) (*StreamSettingResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/stream_setting?index=%s&action=query", index)
	return c.sendStreamRequest(endpoint)
}

// Helper function to send requests related to stream settings
func (c *CameraClient) sendStreamRequest(endpoint string) (*StreamSettingResponse, error) {
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var response StreamSettingResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
