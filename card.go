package zcam

import (
	"fmt"
)

type cardManagementResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  int    `json:"msg"`
}

// CheckCardPresence checks if a storage card is present in the camera
func (c *CameraClient) CheckCardPresence() (bool, error) {
	r, err := c.sendCardRequest("/ctrl/card?action=present")
	if err != nil {
		return false, err
	}

	return r.Code == 0, nil
}

// FormatCard formats the storage card based on its capacity
func (c *CameraClient) FormatCard() (*cardManagementResponse, error) {
	return c.sendCardRequest("/ctrl/card?action=format")
}

// FormatCardAs formats the card specifically to either 'fat32' or 'exfat'
func (c *CameraClient) FormatCardAs(fileSystem string) error {
	if fileSystem != "fat32" && fileSystem != "exfat" {
		return fmt.Errorf("invalid file system type: %s", fileSystem)
	}
	endpoint := fmt.Sprintf("/ctrl/card?action=%s", fileSystem)
	r, err := c.sendCardRequest(endpoint)
	if err != nil {
		return err
	}

	if r.Code != 0 {
		return fmt.Errorf("unexpected code %d", r.Code)
	}

	return nil
}

// QueryCardFreeSpace queries the free space on the card
func (c *CameraClient) QueryCardFreeSpace() (int, error) {
	return c.queryCardSpace("query_free")

}

// QueryCardTotalSpace queries the total space on the card
func (c *CameraClient) QueryCardTotalSpace() (int, error) {
	return c.queryCardSpace("query_total")
}

func (c *CameraClient) queryCardSpace(action string) (int, error) {
	if action != "query_free" && action != "query_total" {
		return -1, fmt.Errorf("invalid query action: %s", action)
	}
	endpoint := fmt.Sprintf("/ctrl/card?action=%s", action)
	r, err := c.sendCardRequest(endpoint)
	if err != nil {
		return -1, err
	}

	return r.Msg, nil
}

// sendCardRequest sends a GET request to the card management endpoints and parses the response
func (c *CameraClient) sendCardRequest(endpoint string) (*cardManagementResponse, error) {
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var response cardManagementResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
