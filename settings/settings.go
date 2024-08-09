package settings

import "fmt"

// Setting represents a camera setting.
type Setting string

// Video settings
const (
	// MovFmtSetting sets the format (type: choice, options: 4KP30/4KP60/...).
	MovFmtSetting Setting = "movfmt"
	// ResolutionSetting sets the resolution (type: choice, options: 4K/C4K/...).
	ResolutionSetting Setting = "resolution"
	// ProjectFPSSetting sets the project frame rate (type: choice, options: 23.98/24/...).
	ProjectFPSSetting Setting = "project_fps"
	// RecordFileFormatSetting sets the file format for recording (type: choice, options: MOV/MP4).
	RecordFileFormatSetting Setting = "record_file_format"
	// RecProxyFileSetting enables recording of a proxy file (type: choice).
	RecProxyFileSetting Setting = "rec_proxy_file"
	// VideoEncoderSetting sets the video encoder (type: choice, options: h264/h265/...).
	VideoEncoderSetting Setting = "video_encoder"
	// SplitDurationSetting sets the video record split duration (type: choice).
	SplitDurationSetting Setting = "split_duration"
	// BitrateLevelSetting sets the bitrate level (type: choice, options: low/medium/high).
	BitrateLevelSetting Setting = "bitrate_level"
	// ComposeModeSetting sets the compose mode (type: choice, options: Normal/WDR).
	ComposeModeSetting Setting = "compose_mode"
	// MovVFRSetting enables or disables variable framerate (type: choice).
	MovVFRSetting Setting = "movvfr"
	// RecFPSSetting sets the playback framerate (type: choice).
	RecFPSSetting Setting = "rec_fps"
	// VideoTLIntervalSetting sets the video timelapse interval (type: range).
	VideoTLIntervalSetting Setting = "video_tl_interval"
	// EnableVideoTLSetting checks if the camera supports video timelapse (type: choice).
	EnableVideoTLSetting Setting = "enable_video_tl"
	// RecDurationSetting sets the recording duration, in seconds (type: range).
	RecDurationSetting Setting = "rec_duration"
	// LastFileNameSetting queries the last recorded file name (type: string).
	LastFileNameSetting Setting = "last_file_name"
)

// Focus & Zoom settings
const (
	// FocusSetting sets the focus mode (type: choice, options: AF/MF).
	FocusSetting Setting = "focus"
	// AFModeSetting sets the autofocus mode (type: choice, options: Flexible Zone/Human Detection).
	AFModeSetting Setting = "af_mode"
	// MFDriveSetting moves the focus plane far/near (type: range).
	MFDriveSetting Setting = "mf_drive"
	// LensZoomSetting controls the lens zoom in/out (type: choice).
	LensZoomSetting Setting = "lens_zoom"
	// OISModeSetting sets the lens optical image stabilization mode (type: choice).
	OISModeSetting Setting = "ois_mode"
	// AFLockSetting locks/unlocks autofocus (type: choice).
	AFLockSetting Setting = "af_lock"
	// LensZoomPosSetting sets the lens zoom position (type: range).
	LensZoomPosSetting Setting = "lens_zoom_pos"
	// LensFocusPosSetting sets the lens focus position (type: range).
	LensFocusPosSetting Setting = "lens_focus_pos"
	// LensFocusSpdSetting controls the speed of MFDrive/LensFocusPos (type: range).
	LensFocusSpdSetting Setting = "lens_focus_spd"
	// CAFSetting enables or disables continuous autofocus (type: choice).
	CAFSetting Setting = "caf"
	// CAFSensSetting sets the sensitivity of continuous autofocus (type: choice).
	CAFSensSetting Setting = "caf_sens"
	// LiveCAFSetting turns continuous autofocus on or off (type: choice).
	LiveCAFSetting Setting = "live_caf"
	// MFMagSetting magnifies the preview when tuning the manual focus (type: choice).
	MFMagSetting Setting = "mf_mag"
	// RestoreLensPosSetting restores the lens position after reboot (type: choice).
	RestoreLensPosSetting Setting = "restore_lens_pos"
)

// Exposure settings
const (
	// MeterModeSetting sets the automatic exposure meter mode (type: choice).
	MeterModeSetting Setting = "meter_mode"
	// MaxISOSetting sets the maximum ISO value (type: choice).
	MaxISOSetting Setting = "max_iso"
	// EVChoiceSetting sets the exposure value (type: choice, options: -3/.../0/.../3).
	EVChoiceSetting Setting = "ev_choice"
	// ISOSetting sets the ISO mode (type: choice, options: Auto/.../200/.../Max ISO).
	ISOSetting Setting = "iso"
	// IrisSetting sets the aperture size (type: choice).
	IrisSetting Setting = "iris"
	// ShutterAngleSetting sets the shutter angle (type: choice, options: Auto/.../45/90/.../360).
	ShutterAngleSetting Setting = "shutter_angle"
	// MaxExpShutterAngleSetting sets the maximum video shutter angle (type: choice).
	MaxExpShutterAngleSetting Setting = "max_exp_shutter_angle"
	// ShutterTimeSetting sets the shutter time (type: choice).
	ShutterTimeSetting Setting = "shutter_time"
	// MaxExpShutterTimeSetting sets the maximum video shutter time (type: choice).
	MaxExpShutterTimeSetting Setting = "max_exp_shutter_time"
	// ShtOperationSetting selects between speed or angle for shutter operation (type: choice).
	ShtOperationSetting Setting = "sht_operation"
	// DualISOSetting enables or disables dual ISO mode (type: choice, options: Auto/Low/High).
	DualISOSetting Setting = "dual_iso"
	// AEFreezeSetting locks/unlocks automatic exposure (type: choice).
	AEFreezeSetting Setting = "ae_freeze"
	// LiveAEFNoSetting shows the live value of the F-number, read-only (type: string).
	LiveAEFNoSetting Setting = "live_ae_fno"
	// LiveAEISOSetting shows the live value of ISO, read-only (type: string).
	LiveAEISOSetting Setting = "live_ae_iso"
	// LiveAEShutterSetting shows the live value of shutter time, read-only (type: string).
	LiveAEShutterSetting Setting = "live_ae_shutter"
	// LiveAEShutterAngleSetting shows the live value of shutter angle, read-only (type: string).
	LiveAEShutterAngleSetting Setting = "live_ae_shutter_angle"
)

// White Balance settings
const (
	// WBSetting sets the white balance mode (type: choice, options: Auto/Manual).
	WBSetting Setting = "wb"
	// MWBSetting sets the manual white balance in kelvin (type: range).
	MWBSetting Setting = "mwb"
	// TintSetting sets the manual white balance tint (type: range).
	TintSetting Setting = "tint"
	// WBPrioritySetting sets the white balance priority (type: choice, options: Ambiance/White).
	WBPrioritySetting Setting = "wb_priority"
	// MWBRSetting sets the manual white balance red gain (type: range).
	MWBRSetting Setting = "mwb_r"
	// MWBGSetting sets the manual white balance green gain (type: range).
	MWBGSetting Setting = "mwb_g"
	// MWBBSetting sets the manual white balance blue gain (type: range).
	MWBBSetting Setting = "mwb_b"
)

// Image settings
const (
	// SharpnessSetting sets the sharpness level (type: choice, options: Strong/Normal/Weak).
	SharpnessSetting Setting = "sharpness"
	// ContrastSetting sets the contrast level (type: range).
	ContrastSetting Setting = "contrast"
	// SaturationSetting sets the saturation level (type: range).
	SaturationSetting Setting = "saturation"
	// BrightnessSetting sets the brightness level (type: range).
	BrightnessSetting Setting = "brightness"
	// LUTSetting sets the lookup table (type: choice, options: rec709/zlog).
	LUTSetting Setting = "lut"
	// LumaLevelSetting sets the luma level (type: choice, options: 0-255/16-235).
	LumaLevelSetting Setting = "luma_level"
	// VignetteSetting applies a vignette effect (type: choice, not supported in E2).
	VignetteSetting Setting = "vignette"
)

// Stream settings
const (
	// SendStreamSetting selects the stream (type: choice, options: stream0/stream1).
	SendStreamSetting Setting = "send_stream"
)

// Audio settings
const (
	// PrimaryAudioSetting sets the primary audio format (type: choice, options: AAC/PCM).
	PrimaryAudioSetting Setting = "primary_audio"
	// AudioChannelSetting sets the audio input channel (type: choice).
	AudioChannelSetting Setting = "audio_channel"
	// AudioInputGainSetting sets the audio input gain level (type: range).
	AudioInputGainSetting Setting = "audio_input_gain"
	// AudioOutputGainSetting sets the audio output gain level (type: range).
	AudioOutputGainSetting Setting = "audio_output_gain"
	// AudioPhantomPowerSetting turns audio phantom power on or off (type: choice).
	AudioPhantomPowerSetting Setting = "audio_phantom_power"
	// AINGainTypeSetting selects the audio gain type (type: choice, options: AGC/MGC).
	AINGainTypeSetting Setting = "ain_gain_type"
)

// Timecode settings
const (
	// TCCountUpSetting sets the timecode count mode (type: choice, options: free run/record run).
	TCCountUpSetting Setting = "tc_count_up"
	// TCHDMIDisplaySetting displays timecode on HDMI (type: choice).
	TCHDMIDisplaySetting Setting = "tc_hdmi_dispaly"
	// TCDropFrameSetting selects timecode drop frame mode (type: choice, options: DF/NDF).
	TCDropFrameSetting Setting = "tc_drop_frame"
)

// Assist tool settings
const (
	// AssistToolDisplaySetting turns the assist tool display on or off (type: choice).
	AssistToolDisplaySetting Setting = "assitool_display"
	// AssistToolPeakOnOffSetting turns the peaking assist on or off (type: choice).
	AssistToolPeakOnOffSetting Setting = "assitool_peak_onoff"
	// AssistToolPeakColorSetting sets the peaking assist color (type: choice).
	AssistToolPeakColorSetting Setting = "assitool_peak_color"
	// AssistToolExposureSetting sets the exposure assist (type: choice, options: Zebra/False Color).
	AssistToolExposureSetting Setting = "assitool_exposure"
	// AssistToolZebraTH1Setting sets the Zebra high value threshold (type: range).
	AssistToolZebraTH1Setting Setting = "assitool_zera_th1"
	// AssistToolZebraTH2Setting sets the Zebra low value threshold (type: range).
	AssistToolZebraTH2Setting Setting = "assitool_zera_th2"
)

// Misc settings
const (
	// SSIDSetting sets the Wi-Fi SSID (type: string).
	SSIDSetting Setting = "ssid"
	// FlickerSetting sets flicker reduction (type: choice, options: 50Hz/60Hz).
	FlickerSetting Setting = "flicker"
	// VideoSystemSetting sets the video system (type: choice, options: NTSC/PAL/CINEMA).
	VideoSystemSetting Setting = "video_system"
	// WiFiSetting turns Wi-Fi on or off (type: choice).
	WiFiSetting Setting = "wifi"
	// BatterySetting shows the battery percentage (type: range).
	BatterySetting Setting = "battery"
	// BatterySetting shows the battery voltage, the value need to divided by 10
	// to get the value in volts (type: range).
	BatteryVoltage Setting = "battery_voltage"
	// LEDSetting turns the LED on or off (type: choice).
	LEDSetting Setting = "led"
	// LCDBacklightSetting sets the LCD backlight level (type: range).
	LCDBacklightSetting Setting = "lcd_backlight"
	// HDMIFormatSetting sets the HDMI format (type: choice, options: Auto/4KP60/4KP30/...).
	HDMIFormatSetting Setting = "hdmi_fmt"
	// HDMIOSDSetting turns the HDMI on-screen display on or off (type: choice).
	HDMIOSDSetting Setting = "hdmi_osd"
	// USBDeviceRoleSetting sets the USB device role (type: choice, options: Host/Mass storage/Network).
	USBDeviceRoleSetting Setting = "usb_device_role"
	// UARTRoleSetting sets the UART role (type: choice, options: Pelco D/Controller).
	UARTRoleSetting Setting = "uart_role"
	// AutoOffSetting enables or disables camera auto off (type: choice).
	AutoOffSetting Setting = "auto_off"
	// AutoOffLCDSetting enables or disables LCD auto off (type: choice).
	AutoOffLCDSetting Setting = "auto_off_lcd"
	// SerialNumberSetting sets the serial number of the camera (type: string).
	SerialNumberSetting Setting = "sn"
	// DesqueezeSetting sets the desqueeze display ratio (type: choice, options: 1x/1.33x/1.5x/2x).
	DesqueezeSetting Setting = "desqueeze"
)

// Multiple Camera settings
const (
	// MultipleModeSetting sets the multiple camera mode (type: choice, options: single/master/slave).
	MultipleModeSetting Setting = "multiple_mode"
	// MultipleIDSetting sets the multiple camera ID (type: range).
	MultipleIDSetting Setting = "multiple_id"
)

// Photo Settings (not supported in E2)
const (
	// PhotoSizeSetting sets the photo resolution (type: choice).
	PhotoSizeSetting Setting = "photosize"
	// PhotoQualitySetting sets the photo quality (type: choice, options: JPEG/RAW).
	PhotoQualitySetting Setting = "photo_q"
	// BurstSetting sets the burst mode (type: choice).
	BurstSetting Setting = "burst"
	// MaxExposureSetting sets the maximum exposure time (type: choice).
	MaxExposureSetting Setting = "max_exp"
	// ShootModeSetting sets the AE exposure mode (type: choice, options: P/A/S/M).
	ShootModeSetting Setting = "shoot_mode"
	// DriveModeSetting sets the drive mode (type: choice, options: single/burst/timelapse).
	DriveModeSetting Setting = "drive_mode"
	// PhotoTLIntervalSetting sets the photo timelapse interval (type: range).
	PhotoTLIntervalSetting Setting = "photo_tl_interval"
	// PhotoTLNumSetting sets the photo timelapse number (type: range).
	PhotoTLNumSetting Setting = "photo_tl_num"
	// PhotoSelfIntervalSetting sets the interval for selfie (type: range).
	PhotoSelfIntervalSetting Setting = "photo_self_interval"
)

func main() {
	fmt.Println("Camera settings constants of type 'Setting' with descriptions defined.")
}
