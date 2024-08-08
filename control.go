package zcam

import "fmt"

// TriggerAutoFocus initiates autofocus
func (c *CameraClient) TriggerAutoFocus() (*CameraControlResponse, error) {
	return c.sendControlRequest("/ctrl/af")
}

// UpdateAutoFocusROI updates the Region of Interest for autofocus
func (c *CameraClient) UpdateAutoFocusROI(x, y, w, h int) (*CameraControlResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi&x=%d&y=%d&w=%d&h=%d", x, y, w, h)
	return c.sendControlRequest(endpoint)
}

// UpdateAutoFocusCenter updates the center point for autofocus ROI
func (c *CameraClient) UpdateAutoFocusCenter(x, y int) (*CameraControlResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/af?action=update_roi_center&x=%d&y=%d", x, y)
	return c.sendControlRequest(endpoint)
}

// QueryAutoFocusROI queries the current Region of Interest for autofocus
func (c *CameraClient) QueryAutoFocusROI() (*CameraControlResponse, error) {
	return c.sendControlRequest("/ctrl/af?action=query")
}

// SetManualFocusDrive adjusts the manual focus in specified increments
func (c *CameraClient) SetManualFocusDrive(drive int) (*CameraControlResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/set?mf_drive=%d", drive)
	return c.sendControlRequest(endpoint)
}

// SetLensFocusPosition sets the focus plane to a specific position
func (c *CameraClient) SetLensFocusPosition(position int) (*CameraControlResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/set?lens_focus_pos=%d", position)
	return c.sendControlRequest(endpoint)
}

// ZoomControl performs zoom actions such as in, out, or stop
func (c *CameraClient) ZoomControl(action string) (*CameraControlResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom=%s", action)
	return c.sendControlRequest(endpoint)
}

// SetZoomPosition sets the zoom to a specific position within a valid range
func (c *CameraClient) SetZoomPosition(position int) (*CameraControlResponse, error) {
	if position < 0 || position > 31 {
		return nil, fmt.Errorf("zoom position out of range")
	}
	endpoint := fmt.Sprintf("/ctrl/set?lens_zoom_pos=%d", position)
	return c.sendControlRequest(endpoint)
}

// sendControlRequest sends a generic GET request and parses the response
func (c *CameraClient) sendControlRequest(endpoint string) (*CameraControlResponse, error) {
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var response CameraControlResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
