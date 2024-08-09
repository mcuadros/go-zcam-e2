package zcam

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryRemainingRecordingTime(t *testing.T) {
	cli := NewCamera(CameraIP)

	result, err := cli.QueryRemainingRecordingTime(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, -1, result)
}

func TestGetSettingChoice(t *testing.T) {
	cli := NewCamera(CameraIP)

	s, err := cli.GetSetting(context.Background(), "record_file_format")
	assert.NoError(t, err)
	assert.NotEqual(t, -1, s)
	assert.Len(t, s.Options, 2)
	assert.NotEqual(t, "", s.MustValueString())
}

func TestGetSettingString(t *testing.T) {
	cli := NewCamera(CameraIP)

	s, err := cli.GetSetting(context.Background(), "sn")
	assert.NoError(t, err)
	assert.NotEqual(t, -1, s)
	assert.Len(t, s.Options, 0)
	assert.NotEqual(t, "", s.MustValueString())
}

func TestGetSettingRange(t *testing.T) {
	cli := NewCamera(CameraIP)

	s, err := cli.GetSetting(context.Background(), "contrast")
	assert.NoError(t, err)
	assert.NotEqual(t, -1, s)
	assert.Len(t, s.Options, 0)
	assert.NotEqual(t, 0, s.MustValueInt())
	assert.NotEqual(t, 0, s.Max)
}

func TestSetSetting(t *testing.T) {
	cli := NewCamera(CameraIP)

	result, err := cli.GetSetting(context.Background(), "record_file_format")
	assert.NoError(t, err)
	assert.Len(t, result.Options, 2)

	err = cli.SetSetting(context.Background(), "record_file_format", result.Value.(string))
	assert.NoError(t, err)
}

func TestQueryTemperature(t *testing.T) {
	cli := NewCamera(CameraIP)

	result, err := cli.QueryTemperature(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, result, 0)
}
