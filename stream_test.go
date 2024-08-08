package zcam

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryStreamSetting(t *testing.T) {
	cli := NewCameraClient(CameraIP)
	is, err := cli.QueryStreamSetting(Stream1)
	require.NoError(t, err)
	fmt.Println(is)
}
