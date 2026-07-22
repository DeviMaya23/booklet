package bookleaf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/devi/booklet/internal/usecase"
)

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

func (c *Client) GetPublicFolders(ctx context.Context, userID string) (*usecase.FolderList, error) {
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

	var result usecase.FolderList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
