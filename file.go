package zcam

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrFileNotOpen   = errors.New("open file before using it")
	ErrUnknownFormat = errors.New("unknown format")
)

type File struct {
	cli          *CameraClient
	folder, file string
	io.ReadCloser
}

type Format string

const (
	// Original image
	Original Format = ""
	// Thumbnail image
	Thumbnail Format = "thumbnail"
	// Screennail is a larger JPEG than the thumbnail
	Screennail Format = "screennail"
)

func (f *File) Open(format Format) error {
	var err error
	switch format {
	case Original:
		f.ReadCloser, err = f.cli.OpenFile(f.folder, f.file)
	case Thumbnail:
		f.ReadCloser, err = f.cli.OpenThumbnail(f.folder, f.file)
	case Screennail:
		f.ReadCloser, err = f.cli.OpenScreennail(f.folder, f.file)
	default:
		return ErrUnknownFormat
	}

	return err
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.ReadCloser == nil {
		return -1, ErrFileNotOpen
	}

	return f.Read(p)
}

func (f *File) Close() error {
	if f.ReadCloser != nil {
		return f.Close()
	}

	return nil
}

func (f *File) Info() (*FileInformation, error) {
	return f.cli.GetFileInfo(f.folder, f.file)
}

func (f *File) CreatedAt() (time.Time, error) {
	info, err := f.cli.GetFileCreationTime(f.folder, f.file)
	if err != nil {
		return time.Time{}, err
	}

	unix, err := strconv.Atoi(info.Msg)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(unix), 0), nil
}

func (f *File) Delete() error {
	defer f.Close()
	_, err := f.cli.DeleteFile(f.folder, f.file)
	if err != nil {
		return err
	}

	return nil
}

type fileListResponse struct {
	Code  int      `json:"code"`
	Desc  string   `json:"desc"`
	Files []string `json:"files"`
}

// FileInformation for handling detailed file information
type FileInformation struct {
	Code        int    `json:"code"`
	Desc        string `json:"desc"`
	Msg         string `json:"msg"`
	Width       int    `json:"w"`
	Height      int    `json:"h"`
	Timescale   int    `json:"vts"`
	PacketCount int    `json:"vcnt"`
	Duration    int    `json:"dur"`
}

// ListFolders lists the folders in the DCIM directory
func (c *CameraClient) ListFolders() ([]string, error) {
	r, err := c.sendFileRequest("/DCIM/")
	if err != nil {
		return nil, err
	}

	return r.Files, nil
}

// ListFiles lists the files in a specific folder
func (c *CameraClient) ListFiles(folder string) ([]*File, error) {
	endpoint := fmt.Sprintf("/DCIM/%s", folder)
	r, err := c.sendFileRequest(endpoint)
	if err != nil {
		return nil, err
	}

	var files []*File
	for _, f := range r.Files {
		files = append(files, &File{cli: c, folder: folder, file: f})
	}

	return files, err
}

// OpenFile downloads a specific file from a given folder.
func (c *CameraClient) OpenFile(folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s", folder, filename)
	return c.getReader(endpoint)
}

// DeleteFile deletes a specific file from a given folder
func (c *CameraClient) DeleteFile(folder, filename string) (*CameraControlResponse, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=rm", folder, filename)
	return c.sendControlRequest(endpoint)
}

// OpenThumbnail fetches the thumbnail of a video file
func (c *CameraClient) OpenThumbnail(folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=thm", folder, filename)
	return c.getReader(endpoint)
}

// OpenScreennail fetches a larger JPEG (screennail) of a video file
func (c *CameraClient) OpenScreennail(folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=scr", folder, filename)
	return c.getReader(endpoint)
}

// GetFileCreationTime gets the creation time of a video file
func (c *CameraClient) GetFileCreationTime(folder, filename string) (*FileInformation, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=ct", folder, filename)
	return c.sendFileInfoRequest(endpoint)
}

// GetFileInfo fetches the video file information including dimensions and duration
func (c *CameraClient) GetFileInfo(folder, filename string) (*FileInformation, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=info", folder, filename)
	return c.sendFileInfoRequest(endpoint)
}

func (c *CameraClient) getReader(endpoint string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *CameraClient) sendFileRequest(endpoint string) (*fileListResponse, error) {
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var r fileListResponse
	if err := decodeJSON(body, &r); err != nil {
		return nil, err
	}

	if r.Code != 0 {
		return nil, fmt.Errorf("unexpected code response %d", r.Code)
	}

	return &r, nil
}

func (c *CameraClient) sendFileInfoRequest(endpoint string) (*FileInformation, error) {
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var r FileInformation
	if err := decodeJSON(body, &r); err != nil {
		return nil, err
	}

	if r.Code != 0 {
		return nil, fmt.Errorf("unexpected code response %d", r.Code)
	}

	return &r, nil
}
