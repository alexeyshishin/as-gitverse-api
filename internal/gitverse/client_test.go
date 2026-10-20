package gitverse

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newTestClient(fn roundTripFunc) *Client {
	c := NewClient("test-token")
	c.httpClient = &http.Client{Transport: fn}
	return c
}

func jsonResp(status int, body string, headers map[string]string) *http.Response {
	resp := &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{},
	}
	for k, v := range headers {
		resp.Header.Set(k, v)
	}
	return resp
}

func TestDo_Success(t *testing.T) {
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("missing auth header")
		}
		return jsonResp(200, `{"id":1,"name":"foo"}`, nil), nil
	})

	var repo Repository
	_, err := c.do(context.Background(), "GET", "/repos/x/foo", nil, &repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Name != "foo" {
		t.Fatalf("got name %q", repo.Name)
	}
}

func TestDo_NotFound(t *testing.T) {
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		return jsonResp(404, `{"code":"not_found","message":"no such repo"}`, nil), nil
	})

	_, err := c.do(context.Background(), "GET", "/repos/x/missing", nil, nil)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDo_RetriesOn429ThenSucceeds(t *testing.T) {
	calls := 0
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return jsonResp(429, `{"code":"rate_limited"}`, map[string]string{"Retry-After": "0"}), nil
		}
		return jsonResp(200, `{"id":1,"name":"ok"}`, nil), nil
	})

	var repo Repository
	_, err := c.do(context.Background(), "GET", "/repos/x/foo", nil, &repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestDo_RetriesOn500ThenFails(t *testing.T) {
	calls := 0
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return jsonResp(500, `boom`, nil), nil
	})

	_, err := c.do(context.Background(), "GET", "/repos/x/foo", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
	if calls != maxRetries+1 {
		t.Fatalf("expected %d calls, got %d", maxRetries+1, calls)
	}
}

func TestDo_MalformedBody(t *testing.T) {
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		return jsonResp(200, `not json`, nil), nil
	})

	var repo Repository
	_, err := c.do(context.Background(), "GET", "/repos/x/foo", nil, &repo)
	if err == nil {
		t.Fatalf("expected decode error")
	}
}
