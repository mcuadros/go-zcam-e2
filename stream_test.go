package zcam

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryStreamSetting(t *testing.T) {
	cli := NewCamera(CameraIP)
	is, err := cli.QueryStreamSetting(context.Background(), Stream1)
	require.NoError(t, err)
	fmt.Println(is)
}
