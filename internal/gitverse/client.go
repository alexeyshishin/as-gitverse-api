package gitverse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://api.gitverse.ru"
	acceptHeader   = "application/vnd.gitverse.object+json;version=1"
	maxRetries     = 3
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	verbose    bool
}

func NewClient(token string) *Client {
	return &Client{
		baseURL: defaultBaseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SetVerbose(verbose bool) {
	c.verbose = verbose
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) (http.Header, error) {
	return c.doURL(ctx, method, c.baseURL+path, body, out)
}

func (c *Client) doURL(ctx context.Context, method, url string, body, out any) (http.Header, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("gitverse: encode request body: %w", err)
		}
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoffDelay(attempt)):
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, fmt.Errorf("gitverse: build request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Accept", acceptHeader)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("gitverse: request failed: %w", err)
			continue
		}

		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = fmt.Errorf("gitverse: read response: %w", readErr)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = (&APIError{HTTPStatus: resp.StatusCode, Code: "rate_limited", Message: string(respBody)}).Classify()
			if attempt < maxRetries {
				wait := retryAfterDelay(resp.Header.Get("Retry-After"), attempt)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(wait):
				}
			}
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("gitverse: server error %d: %s", resp.StatusCode, string(respBody))
			continue
		}

		if resp.StatusCode >= 400 {
			apiErr := &APIError{HTTPStatus: resp.StatusCode}
			if jsonErr := json.Unmarshal(respBody, apiErr); jsonErr != nil {
				apiErr.Message = describeNonJSONBody(respBody, c.verbose)
			}
			return nil, apiErr.Classify()
		}

		if out != nil && len(respBody) > 0 {
			if err := json.Unmarshal(respBody, out); err != nil {
				return nil, fmt.Errorf("gitverse: decode response: %w", err)
			}
		}
		return resp.Header, nil
	}

	return nil, lastErr
}

func describeNonJSONBody(body []byte, verbose bool) string {
	kind := "non-JSON response"
	if looksLikeHTML(body) {
		kind = "HTML response instead of JSON (likely a wrong URL or endpoint)"
	}
	if verbose {
		return fmt.Sprintf("%s: %s", kind, string(body))
	}
	return fmt.Sprintf("%s (re-run with --verbose to see the raw body)", kind)
}

func looksLikeHTML(body []byte) bool {
	trimmed := bytes.ToLower(bytes.TrimSpace(body))
	return bytes.HasPrefix(trimmed, []byte("<!doctype")) || bytes.HasPrefix(trimmed, []byte("<html"))
}

func backoffDelay(attempt int) time.Duration {
	return time.Duration(attempt) * 500 * time.Millisecond
}

func retryAfterDelay(header string, attempt int) time.Duration {
	if header != "" {
		if secs, err := strconv.Atoi(header); err == nil {
			return time.Duration(secs) * time.Second
		}
	}
	return backoffDelay(attempt + 1)
}
