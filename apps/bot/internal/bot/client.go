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
	baseURL              string
	internalHeaderName   string
	internalServiceToken string
	client               *http.Client
}

type HealthcheckResult struct {
	Status string `json:"status"`
}

func NewHTTPMarkerClient(baseURL, internalHeaderName, internalServiceToken string, client *http.Client) *HTTPMarkerClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPMarkerClient{
		baseURL:              baseURL,
		internalHeaderName:   internalHeaderName,
		internalServiceToken: internalServiceToken,
		client:               client,
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

type RuntimeTitleFormatResult struct {
	Format string `json:"format"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("request failed: %d", e.StatusCode)
}

func (c *HTTPMarkerClient) SyncSession(ctx context.Context, channel string) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPost, "/internal/channels/"+channel+"/sessions/sync", nil)
}

func (c *HTTPMarkerClient) GetRuntime(ctx context.Context, channel string) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodGet, "/internal/channels/"+channel+"/runtime", nil)
}

func (c *HTTPMarkerClient) ApplyCurrentTitle(ctx context.Context, channel string) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPost, "/internal/channels/"+channel+"/title/apply-current", nil)
}

func (c *HTTPMarkerClient) RestoreTitle(ctx context.Context, channel string) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPost, "/internal/channels/"+channel+"/title/restore", nil)
}

func (c *HTTPMarkerClient) ToggleTitles(ctx context.Context, channel string) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPost, "/internal/channels/"+channel+"/title/toggle", nil)
}

func (c *HTTPMarkerClient) SetTitleFormat(ctx context.Context, channel, format string, reset bool) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPut, "/internal/channels/"+channel+"/title/format", map[string]any{
		"format": format,
		"reset":  reset,
	})
}

func (c *HTTPMarkerClient) GetTitleFormat(ctx context.Context, channel string) (RuntimeTitleFormatResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/internal/channels/"+channel+"/title/format", nil)
	if err != nil {
		return RuntimeTitleFormatResult{}, fmt.Errorf("build request: %w", err)
	}
	c.applyInternalAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return RuntimeTitleFormatResult{}, fmt.Errorf("title format request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return RuntimeTitleFormatResult{}, decodeAPIError(resp)
	}
	var result RuntimeTitleFormatResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RuntimeTitleFormatResult{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *HTTPMarkerClient) AdvanceSegment(ctx context.Context, channel string) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPost, "/internal/channels/"+channel+"/segments/advance", nil)
}

func (c *HTTPMarkerClient) AddTimelineMarker(ctx context.Context, channel, label string, end bool) (stream.RuntimeState, error) {
	return c.runtimeRequest(ctx, http.MethodPost, "/internal/channels/"+channel+"/markers", map[string]any{
		"label": label,
		"end":   end,
	})
}

func (c *HTTPMarkerClient) Healthcheck(ctx context.Context) (HealthcheckResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return HealthcheckResult{}, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return HealthcheckResult{}, fmt.Errorf("healthcheck request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthcheckResult{}, fmt.Errorf("healthcheck failed: %s", resp.Status)
	}

	var result HealthcheckResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return HealthcheckResult{}, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (c *HTTPMarkerClient) runtimeRequest(ctx context.Context, method, path string, body any) (stream.RuntimeState, error) {
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			return stream.RuntimeState{}, fmt.Errorf("marshal request: %w", err)
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return stream.RuntimeState{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.applyInternalAuth(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return stream.RuntimeState{}, fmt.Errorf("runtime request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return stream.RuntimeState{}, decodeAPIError(resp)
	}
	var result stream.RuntimeState
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return stream.RuntimeState{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

func (c *HTTPMarkerClient) applyInternalAuth(req *http.Request) {
	if c.internalHeaderName != "" && c.internalServiceToken != "" {
		req.Header.Set(c.internalHeaderName, c.internalServiceToken)
	}
}

func decodeAPIError(resp *http.Response) error {
	var payload struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err == nil && payload.Error != "" {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    payload.Error,
		}
	}
	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    resp.Status,
	}
}
