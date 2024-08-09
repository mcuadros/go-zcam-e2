package zcam

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckCardPresence(t *testing.T) {
	cli := NewCameraClient(CameraIP)
	is, err := cli.CheckCardPresence(context.Background())
	require.NoError(t, err)
	require.True(t, is)
}

func TestQueryCardTotalSpace(t *testing.T) {
	cli := NewCameraClient(CameraIP)
	space, err := cli.QueryCardTotalSpace(context.Background())
	require.NoError(t, err)
	require.NotEqual(t, 0, space)
}
