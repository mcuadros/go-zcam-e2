package zcam

import (
	"context"
	"fmt"
	"strconv"
	"time"
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
func (c *Camera) StartVideoRecord() error {
	body, err := c.get("/ctrl/rec?action=start")
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// StopVideoRecord stops video recording or video timelapse recording
func (c *Camera) StopVideoRecord() error {
	body, err := c.get("/ctrl/rec?action=stop")
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// VideoRecord records a video of the give duration, it list the files before
// start to recording and after stop the recording and returns the result.
func (c *Camera) VideoRecord(ctx context.Context, d time.Duration) (*File, error) {
	before, err := c.ListAllFiles()
	if err != nil {
		return nil, fmt.Errorf("unable to list files: %w", err)
	}

	select {
	case <-ctx.Done():
		if err := c.StopVideoRecord(); err != nil {
			return nil, fmt.Errorf("error stopping video, at cacelled context: %w", err)
		}

		return nil, fmt.Errorf("context cancelled")
	case <-time.After(d):
		if err := c.StopVideoRecord(); err != nil {
			return nil, fmt.Errorf("error stopping video: %w", err)
		}
	}

	after, err := c.ListAllFiles()
	if err != nil {
		return nil, fmt.Errorf("unable to list files (after record): %w", err)
	}

	f := compareFileList(before, after)
	if f == nil {
		return nil, fmt.Errorf("record finished but no new files where found")
	}

	return f, nil
}

func compareFileList(a, b []*File) *File {
	seen := map[string]bool{}
	for _, f := range a {
		seen[f.Folder()+f.Filename()] = true
	}

	for _, f := range b {
		if !seen[f.Folder()+f.Filename()] {
			return f
		}
	}

	return nil
}

// QueryRemainingRecordingTime queries the remaining recording time in minutes
func (c *Camera) QueryRemainingRecordingTime() (int, error) {
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
func (c *Camera) GetSetting(key string) (*CameraSetting, error) {
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
func (c *Camera) SetSetting(key, value string) error {
	endpoint := fmt.Sprintf("/ctrl/set?%s=%s", key, value)
	body, err := c.get(endpoint)
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// TriggerAutoFocus initiates autofocus
func (c *Camera) TriggerAutoFocus() error {
	return c.sendControlRequest("/ctrl/af")
}

// UpdateAutoFocusROI updates the Region of Interest for autofocus
func (c *Camera) UpdateAutoFocusROI(x, y, w, h int) error {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi&x=%d&y=%d&w=%d&h=%d", x, y, w, h)
	return c.sendControlRequest(endpoint)
}

// UpdateAutoFocusCenter updates the center point for autofocus ROI
func (c *Camera) UpdateAutoFocusCenter(x, y int) error {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi_center&x=%d&y=%d", x, y)
	return c.sendControlRequest(endpoint)
}

// QueryAutoFocusROI queries the current Region of Interest for autofocus
func (c *Camera) QueryAutoFocusROI() error {
	return c.sendControlRequest("/ctrl/af?action=query")
}

// SetManualFocusDrive adjusts the manual focus in specified increments
func (c *Camera) SetManualFocusDrive(drive int) error {
	endpoint := fmt.Sprintf("/ctrl/set?mf_drive=%d", drive)
	return c.sendControlRequest(endpoint)
}

// SetLensFocusPosition sets the focus plane to a specific position
func (c *Camera) SetLensFocusPosition(position int) error {
	endpoint := fmt.Sprintf("/ctrl/set?lens_focus_pos=%d", position)
	return c.sendControlRequest(endpoint)
}

// ZoomControl performs zoom actions such as in, out, or stop
func (c *Camera) ZoomControl(action string) error {
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom=%s", action)
	return c.sendControlRequest(endpoint)
}

// SetZoomPosition sets the zoom to a specific position within a valid range
func (c *Camera) SetZoomPosition(position int) error {
	if position < 0 || position > 31 {
		return fmt.Errorf("zoom position out of range")
	}
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom_pos=%d", position)
	return c.sendControlRequest(endpoint)
}

// sendControlRequest sends a generic GET request and parses the response
func (c *Camera) sendControlRequest(endpoint string) error {
	body, err := c.get(endpoint)
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}
