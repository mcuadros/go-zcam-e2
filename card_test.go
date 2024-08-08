package zcam

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckCardPresence(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))
	is, err := cli.CheckCardPresence()
	require.NoError(t, err)
	require.True(t, is)
}

func TestQueryCardTotalSpace(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))
	space, err := cli.QueryCardTotalSpace()
	require.NoError(t, err)
	require.NotEqual(t, 0, space)
}
