package gitverse

import (
	"context"
	"net/http"
	"testing"
)

func TestListOrgRepos_BuildsEscapedPath(t *testing.T) {
	var gotPath string
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.EscapedPath()
		return jsonResp(200, `[{"name":"a"}]`, nil), nil
	})

	repos, err := c.ListOrgRepos(context.Background(), "my org")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 || repos[0].Name != "a" {
		t.Fatalf("got %+v", repos)
	}
	if gotPath != "/orgs/my%20org/repos" {
		t.Fatalf("got path %q", gotPath)
	}
}
