package zcam

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFileNotOpen   = errors.New("open file before using it")
	ErrUnknownFormat = errors.New("unknown format")
)

const RootFolder = "/DCIM/"

type Format string

const (
	// Original image
	Original Format = ""
	// Thumbnail image
	Thumbnail Format = "thumbnail"
	// Screennail is a larger JPEG than the thumbnail
	Screennail Format = "screennail"
)

type File struct {
	cli          *Camera
	folder, file string
	io.ReadCloser
}

// NewFileFromValueSetting returns a new File from a SettingValue, usually the
// settings.LastFileNameSetting setting.
func NewFileFromValueSetting(c *Camera, v *SettingValue) (*File, error) {
	return NewFile(c, v.MustValueString())
}

// NewFile returns a new file for a given path.
func NewFile(c *Camera, path string) (*File, error) {
	if strings.Index(path, RootFolder) != 0 {
		return nil, fmt.Errorf("unknown root folder for file %s", path)
	}

	path = path[len(RootFolder):]
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected filename %s", path)
	}

	return &File{cli: c, folder: parts[0], file: parts[1]}, nil
}

func (f *File) Folder() string {
	return f.folder
}

func (f *File) Filename() string {
	return f.file
}

func (f *File) Open(ctx context.Context, format Format) error {
	var err error
	switch format {
	case Original:
		f.ReadCloser, err = f.cli.OpenFile(ctx, f.folder, f.file)
	case Thumbnail:
		f.ReadCloser, err = f.cli.OpenThumbnail(ctx, f.folder, f.file)
	case Screennail:
		f.ReadCloser, err = f.cli.OpenScreennail(ctx, f.folder, f.file)
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

func (f *File) Info(ctx context.Context) (*FileInformation, error) {
	return f.cli.GetFileInfo(ctx, f.folder, f.file)
}

func (f *File) CreatedAt(ctx context.Context) (time.Time, error) {
	info, err := f.cli.GetFileCreationTime(ctx, f.folder, f.file)
	if err != nil {
		return time.Time{}, err
	}

	unix, err := strconv.Atoi(info.Msg)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(unix), 0), nil
}

func (f *File) Delete(ctx context.Context) error {
	defer f.Close()
	return f.cli.DeleteFile(ctx, f.folder, f.file)
}

// Download copy the camera file to a local file, return the downloaded bytes.
// It opens and close this File,
func (f *File) Download(ctx context.Context, format Format, filename string) (int64, error) {
	if err := f.Open(ctx, format); err != nil {
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
func (c *Camera) ListFolders(ctx context.Context) ([]string, error) {
	r, err := c.sendFileRequest(ctx, RootFolder)
	if err != nil {
		return nil, err
	}

	return r.Files, nil
}

// ListFiles lists the files in a specific folder
func (c *Camera) ListFiles(ctx context.Context, folder string) ([]*File, error) {
	endpoint := fmt.Sprintf(RootFolder+"%s", folder)
	r, err := c.sendFileRequest(ctx, endpoint)
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
func (c *Camera) ListAllFiles(ctx context.Context) ([]*File, error) {
	folders, err := c.ListFolders(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving folders: %w", err)
	}

	var files []*File
	for _, folder := range folders {
		f, err := c.ListFiles(ctx, folder)
		if err != nil {
			return nil, fmt.Errorf("error retrieving files from folder %s: %w", folder, err)
		}

		files = append(files, f...)
	}

	return files, nil
}

// OpenFile downloads a specific file from a given folder.
func (c *Camera) OpenFile(ctx context.Context, folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf(RootFolder+"%s/%s", folder, filename)
	return c.getReader(ctx, endpoint)
}

// DeleteFile deletes a specific file from a given folder
func (c *Camera) DeleteFile(ctx context.Context, folder, filename string) error {
	endpoint := fmt.Sprintf(RootFolder+"%s/%s?act=rm", folder, filename)
	return c.sendControlRequest(ctx, endpoint)
}

// OpenThumbnail fetches the thumbnail of a video file
func (c *Camera) OpenThumbnail(ctx context.Context, folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf(RootFolder+"%s/%s?act=thm", folder, filename)
	return c.getReader(ctx, endpoint)
}

// OpenScreennail fetches a larger JPEG (screennail) of a video file
func (c *Camera) OpenScreennail(ctx context.Context, folder, filename string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf(RootFolder+"%s/%s?act=scr", folder, filename)
	return c.getReader(ctx, endpoint)
}

// GetFileCreationTime gets the creation time of a video file
func (c *Camera) GetFileCreationTime(ctx context.Context, folder, filename string) (*FileInformation, error) {
	endpoint := fmt.Sprintf(RootFolder+"%s/%s?act=ct", folder, filename)
	return c.sendFileInfoRequest(ctx, endpoint)
}

// GetFileInfo fetches the video file information including dimensions and duration
func (c *Camera) GetFileInfo(ctx context.Context, folder, filename string) (*FileInformation, error) {
	endpoint := fmt.Sprintf(RootFolder+"%s/%s?act=info", folder, filename)
	return c.sendFileInfoRequest(ctx, endpoint)
}

func (c *Camera) getReader(ctx context.Context, endpoint string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create GET request: %w", err)
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *Camera) sendFileRequest(ctx context.Context, endpoint string) (*fileListResponse, error) {
	body, err := c.get(ctx, endpoint)
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

func (c *Camera) sendFileInfoRequest(ctx context.Context, endpoint string) (*FileInformation, error) {
	body, err := c.get(ctx, endpoint)
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
