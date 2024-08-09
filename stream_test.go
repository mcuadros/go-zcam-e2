package zcam

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryStreamSetting(t *testing.T) {
	cli := NewCamera(CameraIP)
	config, err := cli.QueryStreamSetting(context.Background(), Stream1)
	require.NoError(t, err)
	require.Equal(t, config.Stream, Stream1)
}
