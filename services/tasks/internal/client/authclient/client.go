package authclient

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"singularity.com/pr1/shared/httpx"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type VerifyResponse struct {
	Valid bool `json:"valid"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    httpx.NewClient(3 * time.Second),
	}
}

func (c *Client) Verify(ctx context.Context, token string, requestID string) (bool, int, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+"/v1/auth/verify",
		nil,
	)
	if err != nil {
		return false, 0, err
	}

	req.Header.Set("Authorization", token)
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, resp.StatusCode, nil
	}

	var body VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false, 0, err
	}

	return body.Valid, resp.StatusCode, nil
}
