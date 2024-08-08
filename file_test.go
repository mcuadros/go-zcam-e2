package zcam

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListFolders(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))
	_, err := cli.ListFolders()
	require.NoError(t, err)
}

func TestClientListFiles(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))
	files, err := cli.ListFiles("Z003")
	require.NoError(t, err)
	require.NotEqual(t, files, 0)

	info, err := files[0].Info()
	require.NoError(t, err)
	require.NotEqual(t, info.Height, 0)
}

func TestClienGetFileInfo(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))
	folders, err := cli.GetFileInfo("Z003", "Z003C0002_20240808132826_0001.MOV")
	require.NoError(t, err)
	require.Equal(t, folders.Code, 0)

	fmt.Println(folders)
}
