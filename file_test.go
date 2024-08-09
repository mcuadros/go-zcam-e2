package zcam

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListFolders(t *testing.T) {
	cli := NewCamera(CameraIP)
	_, err := cli.ListFolders(context.Background())
	require.NoError(t, err)
}

func TestClientListFiles(t *testing.T) {
	cli := NewCamera(CameraIP)
	folders, err := cli.ListFolders(context.Background())
	require.NoError(t, err)

	files, err := cli.ListFiles(context.Background(), folders[0])
	require.NoError(t, err)
	require.NotEqual(t, files, 0)

	info, err := files[0].Info(context.Background())
	require.NoError(t, err)
	require.NotEqual(t, info.Height, 0)
}
