package gitverse

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestListCommits_NoQueryParams(t *testing.T) {
	var gotURL *http.Request
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		gotURL = r
		return jsonResp(200, `[{"sha":"abc","author":{"login":"dave"}}]`, nil), nil
	})

	commits, err := c.ListCommits(context.Background(), "owner", "repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 || commits[0].AuthorName() != "dave" {
		t.Fatalf("got %+v", commits)
	}
	if gotURL.URL.Path != "/repos/owner/repo/commits" {
		t.Fatalf("got path %q", gotURL.URL.Path)
	}
	if gotURL.URL.RawQuery != "" {
		t.Fatalf("expected no query params (server-side since/until truncates pagination), got %q", gotURL.URL.RawQuery)
	}
}

func TestListCommits_NotFound(t *testing.T) {
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		return jsonResp(404, `{"code":"not_found"}`, nil), nil
	})

	_, err := c.ListCommits(context.Background(), "owner", "repo")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCommit_AuthorNameFallsBackToCommitAuthor(t *testing.T) {
	var cm Commit
	cm.Commit.Author.Name = "Erin"
	if cm.AuthorName() != "Erin" {
		t.Fatalf("got %q", cm.AuthorName())
	}
}
