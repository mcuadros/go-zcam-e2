package zcam

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	ErrFileNotOpen   = errors.New("open file before using it")
	ErrUnknownFormat = errors.New("unknown format")
)

type File struct {
	cli          *Camera
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

func (f *File) Folder() string {
	return f.folder
}

func (f *File) Filename() string {
	return f.file
}

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

	return f.ReadCloser.Read(p)
}

func (f *File) Close() error {
	if f.ReadCloser != nil {
		return f.ReadCloser.Close()
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
	return f.cli.DeleteFile(f.folder, f.file)
}

// Download copy the camera file to a local file, return the downloaded bytes.
// It opens and close this File,
func (f *File) Download(format Format, filename string) (int64, error) {
	if err := f.Open(format); err != nil {
		return -1, fmt.Errorf("unable to open file %q in folder %q: %w", f.file, f.folder, err)
	}

	defer f.Close()

	file, err := os.Create(filename)
	if err != nil {
		return -1, fmt.Errorf("unable to create file %q: %w", filename, err)
	}

	defer file.Close()

	return io.Copy(file, f)
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
func (c *Camera) ListFolders() ([]string, error) {
	r, err := c.sendFileRequest("/DCIM/")
	if err != nil {
		return nil, err
	}

	return r.Files, nil
}

// ListFiles lists the files in a specific folder
func (c *Camera) ListFiles(folder string) ([]*File, error) {
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

// ListAllFiles lists the files in a all the folders
func (c *Camera) ListAllFiles() ([]*File, error) {
	folders, err := c.ListFolders()
	if err != nil {
		return nil, fmt.Errorf("error retrieving folders: %w", err)
	}

	var files []*File
	for _, folder := range folders {
		f, err := c.ListFiles(folder)
		if err != nil {
			return nil, fmt.Errorf("error retrieving files from folder %s: %w", folder, err)
		}

		files = append(files, f...)
	}

	return files, nil
}

// OpenFile downloads a specific file from a given folder.
func (c *Camera) OpenFile(folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s", folder, filename)
	return c.getReader(endpoint)
}

// DeleteFile deletes a specific file from a given folder
func (c *Camera) DeleteFile(folder, filename string) error {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=rm", folder, filename)
	return c.sendControlRequest(endpoint)
}

// OpenThumbnail fetches the thumbnail of a video file
func (c *Camera) OpenThumbnail(folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=thm", folder, filename)
	return c.getReader(endpoint)
}

// OpenScreennail fetches a larger JPEG (screennail) of a video file
func (c *Camera) OpenScreennail(folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=scr", folder, filename)
	return c.getReader(endpoint)
}

// GetFileCreationTime gets the creation time of a video file
func (c *Camera) GetFileCreationTime(folder, filename string) (*FileInformation, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=ct", folder, filename)
	return c.sendFileInfoRequest(endpoint)
}

// GetFileInfo fetches the video file information including dimensions and duration
func (c *Camera) GetFileInfo(folder, filename string) (*FileInformation, error) {
	endpoint := fmt.Sprintf("/DCIM/%s/%s?act=info", folder, filename)
	return c.sendFileInfoRequest(endpoint)
}

func (c *Camera) getReader(endpoint string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *Camera) sendFileRequest(endpoint string) (*fileListResponse, error) {
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

func (c *Camera) sendFileInfoRequest(endpoint string) (*FileInformation, error) {
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
