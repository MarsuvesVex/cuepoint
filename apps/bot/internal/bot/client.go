package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MarsuvesVex/cuepoint/packages/stream"
)

type HTTPMarkerClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPMarkerClient(baseURL string, client *http.Client) *HTTPMarkerClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPMarkerClient{
		baseURL: baseURL,
		client:  client,
	}
}

func (c *HTTPMarkerClient) CreateMarker(ctx context.Context, input stream.CreateMarkerInput) (CreateMarkerResult, error) {
	body, err := json.Marshal(input)
	if err != nil {
		return CreateMarkerResult{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/markers", bytes.NewReader(body))
	if err != nil {
		return CreateMarkerResult{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return CreateMarkerResult{}, fmt.Errorf("create marker request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return CreateMarkerResult{}, fmt.Errorf("create marker failed: %s", resp.Status)
	}

	var result CreateMarkerResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CreateMarkerResult{}, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}
