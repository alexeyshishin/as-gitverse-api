package gitverse

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestListCollaborators_Success(t *testing.T) {
	var gotPath string
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.EscapedPath()
		return jsonResp(200, `[{"id":1,"login":"bob","permissions":{"pull":true,"push":true,"admin":false}}]`, nil), nil
	})

	collaborators, err := c.ListCollaborators(context.Background(), "owner x", "repo y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collaborators) != 1 || collaborators[0].Login != "bob" || collaborators[0].PermissionLabel() != "write" {
		t.Fatalf("got %+v", collaborators)
	}
	if gotPath != "/repos/owner%20x/repo%20y/collaborators" {
		t.Fatalf("got path %q", gotPath)
	}
}

func TestListCollaborators_NotFound(t *testing.T) {
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		return jsonResp(404, `{"code":"not_found"}`, nil), nil
	})

	_, err := c.ListCollaborators(context.Background(), "owner", "repo")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAddCollaborator_Success(t *testing.T) {
	var gotPath, gotMethod string
	var gotBody map[string]string
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.EscapedPath()
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)
		return jsonResp(204, ``, nil), nil
	})

	err := c.AddCollaborator(context.Background(), "owner", "repo", "charlie", "write")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "PUT" {
		t.Fatalf("got method %q", gotMethod)
	}
	if gotPath != "/repos/owner/repo/collaborators/charlie" {
		t.Fatalf("got path %q", gotPath)
	}
	if gotBody["permission"] != "write" {
		t.Fatalf("got body %+v", gotBody)
	}
}

func TestAddCollaborator_InvalidPermission(t *testing.T) {
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		t.Fatalf("request should not be sent for invalid permission")
		return nil, nil
	})

	err := c.AddCollaborator(context.Background(), "owner", "repo", "charlie", "superadmin")
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
