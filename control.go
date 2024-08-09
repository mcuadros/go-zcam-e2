package zcam

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mcuadros/go-zcam-e2/settings"
)

type ReadOnly bool

func (v *ReadOnly) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "0":
		*v = false
	case "1":
		*v = true
	default:
		return fmt.Errorf("unexpected value %s at ReadOnly type", data)
	}

	return nil
}

type SettingType int

const (
	ChoiceSettingType SettingType = 1
	RangeSettingType  SettingType = 2
	StringSettingType SettingType = 3
)

type SettingValue struct {
	Code     int         `json:"code"`
	Desc     string      `json:"desc"`
	Key      string      `json:"key"`
	Type     SettingType `json:"type"`
	ReadOnly ReadOnly    `json:"ro"`
	Value    any         `json:"value"`
	Options  []string    `json:"opts,omitempty"` // Only for choice type
	Min      int         `json:"min,omitempty"`  // Only for range type
	Max      int         `json:"max,omitempty"`  // Only for range type
	Step     int         `json:"step,omitempty"` // Only for range type
}

func (v *SettingValue) Kind() reflect.Kind {
	return reflect.TypeOf(v.Value).Kind()
}

func (v *SettingValue) MustValueString() string {
	i, ok := v.Value.(string)
	if !ok {
		panic("value it not a string")
	}

	return i
}

func (v *SettingValue) MustValueInt() int {
	i, ok := v.Value.(float64)
	if !ok {
		panic("value it not a int")
	}

	return int(i)
}

func (c *SettingValue) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "Key: %s\n", c.Key)
	fmt.Fprintf(buf, "Read-Only: %t\n", c.ReadOnly)

	switch c.Kind() {
	case reflect.Float64:
		fmt.Fprintf(buf, "Value: %d\n", c.MustValueInt())
	case reflect.String:
		fmt.Fprintf(buf, "Value: %s\n", c.MustValueString())
	}

	switch c.Type {
	case ChoiceSettingType:
		fmt.Fprintf(buf, "Type: choice\n")
		fmt.Fprintf(buf, "Options: \n")
		for _, opt := range c.Options {
			fmt.Fprintf(buf, "  - %s\n", opt)
		}
	case RangeSettingType:
		fmt.Fprintf(buf, "Type: range\n")
		fmt.Fprintf(buf, "Min: %d\n", c.Min)
		fmt.Fprintf(buf, "Max: %d\n", c.Max)
		fmt.Fprintf(buf, "Step: %d\n", c.Step)

	}

	return buf.String()
}

// StartVideoRecord starts video recording or video timelapse recording
func (c *Camera) StartVideoRecord(ctx context.Context) error {
	body, err := c.get(ctx, "/ctrl/rec?action=start")
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// StopVideoRecord stops video recording or video timelapse recording
func (c *Camera) StopVideoRecord(ctx context.Context) error {
	body, err := c.get(ctx, "/ctrl/rec?action=stop")
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

// VideoRecord records a video of the give duration, returns a File from the
// setting value of settings.LastFileNameSetting.
func (c *Camera) VideoRecord(ctx context.Context, d time.Duration) (*File, error) {
	if err := c.StartVideoRecord(ctx); err != nil {
		return nil, fmt.Errorf("error starting video, at cacelled context: %w", err)
	}

	select {
	case <-ctx.Done():
		if err := c.StopVideoRecord(context.Background()); err != nil {
			return nil, fmt.Errorf("error stopping video, at cacelled context: %w", err)
		}

		return nil, fmt.Errorf("context cancelled")
	case <-time.After(d):
		if err := c.StopVideoRecord(ctx); err != nil {
			return nil, fmt.Errorf("error stopping video: %w", err)
		}
	}

	v, err := c.GetSetting(ctx, settings.LastFileNameSetting)
	if err != nil {
		return nil, fmt.Errorf("unable to recover %s setting: %w", settings.LastFileNameSetting, err)
	}

	return NewFileFromValueSetting(c, v)
}

// QueryRemainingRecordingTime queries the remaining recording time.
func (c *Camera) QueryRemainingRecordingTime(ctx context.Context) (time.Duration, error) {
	body, err := c.get(ctx, "/ctrl/rec?action=remain")
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

	mins, err := strconv.Atoi(r.Msg)
	if err != nil {
		return -1, err
	}

	return time.Minute * time.Duration(mins), nil
}

// GetSetting retrieves a camera setting based on its key
func (c *Camera) GetSetting(ctx context.Context, key settings.Setting) (*SettingValue, error) {
	endpoint := fmt.Sprintf("/ctrl/get?k=%s", key)
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var setting SettingValue
	if err := decodeJSON(body, &setting); err != nil {
		return nil, err
	}

	return &setting, nil
}

// SetSetting changes a camera setting for a given key to a specified value
func (c *Camera) SetSetting(ctx context.Context, key settings.Setting, value any) error {
	endpoint := fmt.Sprintf("/ctrl/set?%s=%s", key, convertToString(value))
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}

func convertToString(input any) string {
	switch v := input.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	default:
		return "Unsupported type"
	}
}

// TriggerAutoFocus initiates autofocus
func (c *Camera) TriggerAutoFocus(ctx context.Context) error {
	return c.sendControlRequest(ctx, "/ctrl/af")
}

// UpdateAutoFocusROI updates the Region of Interest for autofocus
func (c *Camera) UpdateAutoFocusROI(ctx context.Context, x, y, w, h int) error {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi&x=%d&y=%d&w=%d&h=%d", x, y, w, h)
	return c.sendControlRequest(ctx, endpoint)
}

// UpdateAutoFocusCenter updates the center point for autofocus ROI
func (c *Camera) UpdateAutoFocusCenter(ctx context.Context, x, y int) error {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi_center&x=%d&y=%d", x, y)
	return c.sendControlRequest(ctx, endpoint)
}

// QueryAutoFocusROI queries the current Region of Interest for autofocus
func (c *Camera) QueryAutoFocusROI(ctx context.Context) error {
	return c.sendControlRequest(ctx, "/ctrl/af?action=query")
}

// SetManualFocusDrive adjusts the manual focus in specified increments
func (c *Camera) SetManualFocusDrive(ctx context.Context, drive int) error {
	endpoint := fmt.Sprintf("/ctrl/set?mf_drive=%d", drive)
	return c.sendControlRequest(ctx, endpoint)
}

// SetLensFocusPosition sets the focus plane to a specific position
func (c *Camera) SetLensFocusPosition(ctx context.Context, position int) error {
	endpoint := fmt.Sprintf("/ctrl/set?lens_focus_pos=%d", position)
	return c.sendControlRequest(ctx, endpoint)
}

// ZoomControl performs zoom actions such as in, out, or stop
func (c *Camera) ZoomControl(ctx context.Context, action string) error {
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom=%s", action)
	return c.sendControlRequest(ctx, endpoint)
}

// SetZoomPosition sets the zoom to a specific position within a valid range
func (c *Camera) SetZoomPosition(ctx context.Context, position int) error {
	if position < 0 || position > 31 {
		return fmt.Errorf("zoom position out of range")
	}
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom_pos=%d", position)
	return c.sendControlRequest(ctx, endpoint)
}

// sendControlRequest sends a generic GET request and parses the response
func (c *Camera) sendControlRequest(ctx context.Context, endpoint string) error {
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return err
	}

	return decodeBasicRequest(body)
}
