package zcam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// CameraClient is a struct that contains the base URL of the camera and an HTTP client
type CameraClient struct {
	BaseURL string
	Client  *http.Client
}

// NewCameraClient initializes and returns a CameraClient
func NewCameraClient(baseURL string) *CameraClient {
	return &CameraClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// get performs a GET request to the given endpoint and returns the response body or an error
func (c *CameraClient) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// decodeJSON is a helper function to decode a JSON response into the provided interface
func decodeJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("error decoding JSON response: %w", err)
	}
	return nil
}

// HealthCheck performs a health check on the camera's API service
func (c *CameraClient) HealthCheck() (*HealthCheckResponse, error) {
	body, err := c.get("/url")
	if err != nil {
		return nil, err
	}

	var result HealthCheckResponse
	if err := decodeJSON(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCameraInfo retrieves and returns the camera information
func (c *CameraClient) GetCameraInfo() (*CameraInfo, error) {
	body, err := c.get("/info")
	if err != nil {
		return nil, err
	}

	var info CameraInfo
	if err := decodeJSON(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// StartSession starts a control session with the camera
func (c *CameraClient) StartSession() (string, error) {
	body, err := c.get("/ctrl/session")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// QuitSession ends the control session with the camera
func (c *CameraClient) QuitSession() (string, error) {
	body, err := c.get("/ctrl/session?action=quit")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// SyncDateTime synchronizes the camera's date and time with the current system time
func (c *CameraClient) SyncDateTime(dateTime time.Time) (string, error) {
	endpoint := fmt.Sprintf("/datetime?date=%s&time=%s", dateTime.Format("2006-01-02"), dateTime.Format("15:04:05"))
	body, err := c.get(endpoint)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ShutdownSystem sends a shutdown command to the camera
func (c *CameraClient) ShutdownSystem() (string, error) {
	body, err := c.get("/ctrl/shutdown")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// RebootSystem sends a reboot command to the camera
func (c *CameraClient) RebootSystem() (string, error) {
	body, err := c.get("/ctrl/reboot")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ChangeWorkingMode switches the camera's working mode based on the provided constant
func (c *CameraClient) ChangeWorkingMode(mode string) (string, error) {
	var endpoint string

	switch mode {
	case ModeVideoRecord:
		endpoint = "/ctrl/mode?action=to_rec"
	case ModePlayback:
		endpoint = "/ctrl/mode?action=to_pb"
	case ModeStandby:
		endpoint = "/ctrl/mode?action=to_standby"
	case ModeExitStandby:
		endpoint = "/ctrl/mode?action=exit_standby"
	default:
		return "", errors.New("unknown mode")
	}

	body, err := c.get(endpoint)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// StartVideoRecord starts video recording or video timelapse recording
func (c *CameraClient) StartVideoRecord() (string, error) {
	body, err := c.get("/ctrl/rec?action=start")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// StopVideoRecord stops video recording or video timelapse recording
func (c *CameraClient) StopVideoRecord() (string, error) {
	body, err := c.get("/ctrl/rec?action=stop")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// QueryRemainingRecordingTime queries the remaining recording time in minutes
func (c *CameraClient) QueryRemainingRecordingTime() (*RemainingTimeResponse, error) {
	body, err := c.get("/ctrl/rec?action=remain")
	if err != nil {
		return nil, err
	}

	var remain RemainingTimeResponse
	if err := decodeJSON(body, &remain); err != nil {
		return nil, err
	}
	return &remain, nil
}

// GetSetting retrieves a camera setting based on its key
func (c *CameraClient) GetSetting(key string) (*CameraSettingResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/get?k=%s", key)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var setting CameraSettingResponse
	if err := decodeJSON(body, &setting); err != nil {
		return nil, err
	}
	return &setting, nil
}

// SetSetting changes a camera setting for a given key to a specified value
func (c *CameraClient) SetSetting(key, value string) (*SetSettingResponse, error) {
	endpoint := fmt.Sprintf("/ctrl/set?%s=%s", key, value)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var response SetSettingResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// SetNetworkMode sets the camera's network mode, using net.IP and net.IPMask
func (c *CameraClient) SetNetworkMode(mode NetworkMode, ipaddr net.IP, netmask net.IPMask, gateway net.IP) (*NetworkInfoResponse, error) {
	var endpoint string
	switch mode {
	case NetworkModeRouter, NetworkModeDirect:
		endpoint = fmt.Sprintf("/ctrl/network?action=set&mode=%s", mode)
	case NetworkModeStatic:
		if ipaddr == nil || netmask == nil || gateway == nil {
			return nil, fmt.Errorf("invalid static network configuration")
		}
		// Convert netmask to CIDR prefix length for compatibility
		prefixLength, _ := netmask.Size()
		endpoint = fmt.Sprintf("/ctrl/network?action=set&mode=%s&ipaddr=%s&netmask=%d&gateway=%s",
			mode, ipaddr.String(), prefixLength, gateway.String())
	default:
		return nil, fmt.Errorf("invalid network mode")
	}

	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var response NetworkInfoResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
