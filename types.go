package zcam

// Camera mode constants
const (
	ModeVideoRecord = "rec"
	ModePlayback    = "pb"
	ModeStandby     = "standby"
	ModeExitStandby = "exit_standby"
)

// HealthCheckResponse struct models the response from the health check endpoint
type HealthCheckResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}

// CameraInfo struct models the camera information returned from the /info endpoint
type CameraInfo struct {
	Model  string `json:"model"`
	Number string `json:"number"`
	Sw     string `json:"sw"`
	Hw     string `json:"hw"`
	Mac    string `json:"mac"`
	EthIP  string `json:"eth_ip"`
	SN     string `json:"sn"`
}

// WorkingModeResponse struct models the response from querying the working mode
type WorkingModeResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"` // Possible values: rec, rec_ing, rec_paused, cap, pb, etc.
}

// RemainingTimeResponse struct models the response from querying the remaining recording time
type RemainingTimeResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  int    `json:"msg"` // Remaining time in minutes
}

type CameraSettingResponse struct {
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

type SetSettingResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
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

const (
	Stream0 = "stream0"
	Stream1 = "stream1"
)

type StreamSettingResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}

const (
	//StreamSetting stream0/stream1
	StreamSetting = "steam"
	// SettingWidth Video width in pixels",
	SettingWidth = "width"
	// Video height in pixels",
	SettingHeight = "height"
	// Encode bitrate in bps",
	SettingBitrate = "bitrate"
	// Frames per second of the stream",
	SettingFPS = "fps"
	// Video encoder (e.g., h264, h265)",
	SettingVenc = "venc"
	// Bit width of the H.265 encoder",
	SettingBitwidth = "bitwidth"
	// Designates the active stream for network streaming
	SettingSendStream = "send_stream"
)

type CameraControlResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}
