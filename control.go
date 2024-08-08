package zcam

import (
	"fmt"
	"strconv"
)

type CameraSetting struct {
	Code  int      `json:"code"`
	Desc  string   `json:"desc"`
	Key   string   `json:"key"`
	Type  int      `json:"type"`
	RO    int      `json:"ro"`
	Value string   `json:"value"`
	Opts  []string `json:"opts,omitempty"` // Only for choice type
	Min   int      `json:"min,omitempty"`  // Only for range type
	Max   int      `json:"max,omitempty"`  // Only for range type
	Step  int      `json:"step,omitempty"` // Only for range type
}

// StartVideoRecord starts video recording or video timelapse recording
func (c *CameraClient) StartVideoRecord() error {
	body, err := c.get("/ctrl/rec?action=start")
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// StopVideoRecord stops video recording or video timelapse recording
func (c *CameraClient) StopVideoRecord() error {
	body, err := c.get("/ctrl/rec?action=stop")
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// QueryRemainingRecordingTime queries the remaining recording time in minutes
func (c *CameraClient) QueryRemainingRecordingTime() (int, error) {
	body, err := c.get("/ctrl/rec?action=remain")
	if err != nil {
		return -1, err
	}

	var r struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := decodeJSON(body, &r); err != nil {
		return -1, err
	}

	return strconv.Atoi(r.Msg)
}

// GetSetting retrieves a camera setting based on its key
func (c *CameraClient) GetSetting(key string) (*CameraSetting, error) {
	endpoint := fmt.Sprintf("/ctrl/get?k=%s", key)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var setting CameraSetting
	if err := decodeJSON(body, &setting); err != nil {
		return nil, err
	}

	return &setting, nil
}

// SetSetting changes a camera setting for a given key to a specified value
func (c *CameraClient) SetSetting(key, value string) error {
	endpoint := fmt.Sprintf("/ctrl/set?%s=%s", key, value)
	body, err := c.get(endpoint)
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// TriggerAutoFocus initiates autofocus
func (c *CameraClient) TriggerAutoFocus() error {
	return c.sendControlRequest("/ctrl/af")
}

// UpdateAutoFocusROI updates the Region of Interest for autofocus
func (c *CameraClient) UpdateAutoFocusROI(x, y, w, h int) error {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi&x=%d&y=%d&w=%d&h=%d", x, y, w, h)
	return c.sendControlRequest(endpoint)
}

// UpdateAutoFocusCenter updates the center point for autofocus ROI
func (c *CameraClient) UpdateAutoFocusCenter(x, y int) error {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi_center&x=%d&y=%d", x, y)
	return c.sendControlRequest(endpoint)
}

// QueryAutoFocusROI queries the current Region of Interest for autofocus
func (c *CameraClient) QueryAutoFocusROI() error {
	return c.sendControlRequest("/ctrl/af?action=query")
}

// SetManualFocusDrive adjusts the manual focus in specified increments
func (c *CameraClient) SetManualFocusDrive(drive int) error {
	endpoint := fmt.Sprintf("/ctrl/set?mf_drive=%d", drive)
	return c.sendControlRequest(endpoint)
}

// SetLensFocusPosition sets the focus plane to a specific position
func (c *CameraClient) SetLensFocusPosition(position int) error {
	endpoint := fmt.Sprintf("/ctrl/set?lens_focus_pos=%d", position)
	return c.sendControlRequest(endpoint)
}

// ZoomControl performs zoom actions such as in, out, or stop
func (c *CameraClient) ZoomControl(action string) error {
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom=%s", action)
	return c.sendControlRequest(endpoint)
}

// SetZoomPosition sets the zoom to a specific position within a valid range
func (c *CameraClient) SetZoomPosition(position int) error {
	if position < 0 || position > 31 {
		return fmt.Errorf("zoom position out of range")
	}
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom_pos=%d", position)
	return c.sendControlRequest(endpoint)
}

// sendControlRequest sends a generic GET request and parses the response
func (c *CameraClient) sendControlRequest(endpoint string) error {
	body, err := c.get(endpoint)
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}
