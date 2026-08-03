package bookleaf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrUnauthorized     = errors.New("bookleaf: unauthorized")
	ErrUnexpectedStatus = errors.New("bookleaf: unexpected status")
)

type Folder struct {
	FolderID   string `json:"folder_id"`
	Token      string `json:"token"`
	FolderName string `json:"folder_name"`
}

type FolderList struct {
	FolderList []Folder `json:"folder_list"`
}

type Client struct {
	httpClient     *http.Client
	host           string
	internalSecret string
}

func NewClient(host, internalSecret string) *Client {
	return &Client{
		httpClient:     &http.Client{},
		host:           host,
		internalSecret: internalSecret,
	}
}

func (c *Client) DeleteAccount(ctx context.Context, kindeUserID string) error {
	url := fmt.Sprintf("%s/internal/accounts/%s", c.host, kindeUserID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Bookleaf-Internal-Secret", c.internalSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: status %d", ErrUnexpectedStatus, resp.StatusCode)
	}
	return nil
}

func (c *Client) GetPublicFolders(ctx context.Context, userID string) (*FolderList, error) {
	url := fmt.Sprintf("%s/internal/users/%s/public-folders", c.host, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Bookleaf-Internal-Secret", c.internalSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get public folders: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get public folders: unexpected status %d", resp.StatusCode)
	}

	var result FolderList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
