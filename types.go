package zcam

// WorkingModeResponse struct models the response from querying the working mode
type WorkingModeResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"` // Possible values: rec, rec_ing, rec_paused, cap, pb, etc.
}

// NetworkMode defines a type for different network modes
type NetworkMode string

// Constants for network modes
const (
	NetworkModeRouter NetworkMode = "Router"
	NetworkModeDirect NetworkMode = "Direct"
	NetworkModeStatic NetworkMode = "Static"
)

// NetworkInfoResponse and NetworkConfigResponse to parse responses from network queries
type NetworkInfoResponse struct {
	Code    int    `json:"code"`
	Desc    string `json:"desc"`
	Mode    string `json:"mode"`
	IP      string `json:"ip,omitempty"`
	Netmask string `json:"netmask,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

type CameraControlResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}
