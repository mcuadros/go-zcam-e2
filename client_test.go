package zcam

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var CameraIP string

func init() {
	CameraIP = os.Getenv("CAMERA_IP")
}

func TestGetCameraInfo(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))

	result, err := cli.GetCameraInfo()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Model)
	assert.NotEmpty(t, result.Number)
	assert.NotEmpty(t, result.Sw)
	assert.NotEmpty(t, result.Hw)
	assert.NotEmpty(t, result.Mac)
	assert.NotEmpty(t, result.EthIP)
	assert.NotEmpty(t, result.SN)
}

// TestStartSession tests the StartSession method
func TestStartSession(t *testing.T) {
	cli := NewCameraClient(fmt.Sprintf("http://%s", CameraIP))

	err := cli.StartSession()
	assert.NoError(t, err)
}

// TestQuitSession tests the QuitSession method
func TestQuitSession(t *testing.T) {
	server := mockServer()
	defer server.Close()

	client := NewCameraClient(server.URL)

	err := client.QuitSession()
	assert.NoError(t, err)
}

// TestSyncDateTime tests the SyncDateTime method
func TestSyncDateTime(t *testing.T) {
	server := mockServer()
	defer server.Close()

	client := NewCameraClient(server.URL)

	result, err := client.SyncDateTime(time.Now())
	assert.NoError(t, err)
	assert.Equal(t, "Date/Time synced successfully.", result)
}

// TestShutdownSystem tests the ShutdownSystem method
func TestShutdownSystem(t *testing.T) {
	server := mockServer()
	defer server.Close()

	client := NewCameraClient(server.URL)

	result, err := client.ShutdownSystem()
	assert.NoError(t, err)
	assert.Equal(t, "System shutdown successfully.", result)
}

// TestRebootSystem tests the RebootSystem method
func TestRebootSystem(t *testing.T) {
	server := mockServer()
	defer server.Close()

	client := NewCameraClient(server.URL)

	result, err := client.RebootSystem()
	assert.NoError(t, err)
	assert.Equal(t, "System rebooted successfully.", result)
}

// TestErrorHandling tests error handling in the get method
func TestErrorHandling(t *testing.T) {
	client := NewCameraClient("http://invalid-url")

	// Expecting an error when trying to make a request to an invalid URL
	_, err := client.get("/invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error making GET request")
}

// mockServer creates a mock HTTP server for testing
func mockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/url":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 0, "desc": "demo", "msg": "OK"}`))
		case "/info":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"model": "TestModel",
				"number": "1",
				"sw": "0.82",
				"hw": "1",
				"mac": "4e:4:b8:2d:78:db",
				"eth_ip": "192.168.9.81",
				"sn": "329A0010009"
			}`))
		case "/ctrl/session":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Session started successfully."))
		case "/ctrl/session?action=quit":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Session quit successfully."))
		case "/datetime":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Date/Time synced successfully."))
		case "/ctrl/shutdown":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("System shutdown successfully."))
		case "/ctrl/reboot":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("System rebooted successfully."))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}
