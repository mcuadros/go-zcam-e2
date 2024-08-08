package zcam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryRemainingRecordingTime(t *testing.T) {
	cli := NewCameraClient(CameraIP)

	result, err := cli.QueryRemainingRecordingTime()
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetSetting(t *testing.T) {
	cli := NewCameraClient(CameraIP)

	result, err := cli.GetSetting("record_file_format")
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
	assert.Len(t, result.Opts, 2)
}

func TestSetSetting(t *testing.T) {
	cli := NewCameraClient(CameraIP)

	result, err := cli.GetSetting("record_file_format")
	assert.NoError(t, err)
	assert.Len(t, result.Opts, 2)

	err = cli.SetSetting("record_file_format", result.Value)
	assert.NoError(t, err)
}

func TestTriggerAutoFocus(t *testing.T) {
	cli := NewCameraClient(CameraIP)

	err := cli.TriggerAutoFocus()
	assert.NoError(t, err)
}
