package zcam

import (
	"context"
	"fmt"
	"strconv"
)

type cardManagementResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Msg  string `json:"msg"`
}

// CheckCardPresence checks if a storage card is present in the camera
func (c *Camera) CheckCardPresence(ctx context.Context) (bool, error) {
	r, err := c.sendCardRequest(ctx, "/ctrl/card?action=present")
	if err != nil {
		return false, err
	}

	return r.Code == 0, nil
}

// FormatCard formats the storage card based on its capacity
func (c *Camera) FormatCard(ctx context.Context) error {
	r, err := c.sendCardRequest(ctx, "/ctrl/card?action=format")
	if err != nil {
		return err
	}

	if r.Code != 0 {
		return fmt.Errorf("unexpected code %d", r.Code)
	}

	return nil
}

// FormatCardAs formats the card specifically to either 'fat32' or 'exfat'
func (c *Camera) FormatCardAs(ctx context.Context, fileSystem string) error {
	if fileSystem != "fat32" && fileSystem != "exfat" {
		return fmt.Errorf("invalid file system type: %s", fileSystem)
	}

	endpoint := fmt.Sprintf("/ctrl/card?action=%s", fileSystem)
	r, err := c.sendCardRequest(ctx, endpoint)
	if err != nil {
		return err
	}

	if r.Code != 0 {
		return fmt.Errorf("unexpected code %d", r.Code)
	}

	return nil
}

// QueryCardFreeSpace queries the free space on the card
func (c *Camera) QueryCardFreeSpace(ctx context.Context) (int, error) {
	return c.queryCardSpace(ctx, "query_free")

}

// QueryCardTotalSpace queries the total space on the card
func (c *Camera) QueryCardTotalSpace(ctx context.Context) (int, error) {
	return c.queryCardSpace(ctx, "query_total")
}

func (c *Camera) queryCardSpace(ctx context.Context, action string) (int, error) {
	if action != "query_free" && action != "query_total" {
		return -1, fmt.Errorf("invalid query action: %s", action)
	}
	endpoint := fmt.Sprintf("/ctrl/card?action=%s", action)
	r, err := c.sendCardRequest(ctx, endpoint)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(r.Msg)
}

// sendCardRequest sends a GET request to the card management endpoints and parses the response
func (c *Camera) sendCardRequest(ctx context.Context, endpoint string) (*cardManagementResponse, error) {
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var response cardManagementResponse
	if err := decodeJSON(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
