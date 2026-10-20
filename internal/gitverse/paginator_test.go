package gitverse

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestParseLinkHeader(t *testing.T) {
	h := http.Header{}
	h.Set("Link", `<https://api.gitverse.ru/user/repos?page=2>; rel="next", <https://api.gitverse.ru/user/repos?page=5>; rel="last"`)

	links := parseLinkHeader(h)
	if links["next"] != "https://api.gitverse.ru/user/repos?page=2" {
		t.Fatalf("got next=%q", links["next"])
	}
	if links["last"] != "https://api.gitverse.ru/user/repos?page=5" {
		t.Fatalf("got last=%q", links["last"])
	}
}

func TestNextPageURL_NoNext(t *testing.T) {
	h := http.Header{}
	if _, ok := nextPageURL("https://api.gitverse.ru/x", h); ok {
		t.Fatalf("expected no next page")
	}
}

func TestNextPageURL_IgnoresLinkHost(t *testing.T) {
	h := http.Header{}
	h.Set("Link", `<https://gitverse.ru/public/api/orgs/acme/repos?page=2>; rel="next"`)

	next, ok := nextPageURL("https://api.gitverse.ru/orgs/acme/repos", h)
	if !ok {
		t.Fatalf("expected next page")
	}
	if next != "https://api.gitverse.ru/orgs/acme/repos?page=2" {
		t.Fatalf("got next=%q", next)
	}
}

func TestFetchAllPages_MultiPage(t *testing.T) {
	page := 1
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		switch page {
		case 1:
			page++
			return jsonResp(200, `[{"name":"a"}]`, map[string]string{
				"Link": `<https://api.gitverse.ru/x?page=2>; rel="next"`,
			}), nil
		case 2:
			page++
			return jsonResp(200, `[{"name":"b"}]`, map[string]string{
				"Link": `<https://api.gitverse.ru/x?page=3>; rel="next"`,
			}), nil
		default:
			return jsonResp(200, `[{"name":"c"}]`, nil), nil
		}
	})

	items, err := fetchAllPages[Repository](context.Background(), c, "https://api.gitverse.ru/x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 || items[0].Name != "a" || items[1].Name != "b" || items[2].Name != "c" {
		t.Fatalf("got %+v", items)
	}
}

func TestFetchAllPages_StopsOnError(t *testing.T) {
	calls := 0
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return jsonResp(404, `{"code":"not_found","message":"missing"}`, nil), nil
	})

	_, err := fetchAllPages[Repository](context.Background(), c, "https://api.gitverse.ru/x")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestFetchAllPages_RespectsMaxPages(t *testing.T) {
	calls := 0
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return jsonResp(200, `[{"name":"x"}]`, map[string]string{
			"Link": `<https://api.gitverse.ru/x?page=next>; rel="next"`,
		}), nil
	})

	items, err := fetchAllPages[Repository](context.Background(), c, "https://api.gitverse.ru/x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != maxPages {
		t.Fatalf("expected %d calls, got %d", maxPages, calls)
	}
	if len(items) != maxPages {
		t.Fatalf("expected %d items, got %d", maxPages, len(items))
	}
}

func TestListRepos_FollowsPagination(t *testing.T) {
	page := 1
	c := newTestClient(func(r *http.Request) (*http.Response, error) {
		switch page {
		case 1:
			page++
			return jsonResp(200, `[{"name":"a"}]`, map[string]string{
				"Link": `<https://api.gitverse.ru/user/repos?page=2>; rel="next"`,
			}), nil
		default:
			return jsonResp(200, `[{"name":"b"}]`, nil), nil
		}
	})

	repos, err := c.ListRepos(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 2 || repos[0].Name != "a" || repos[1].Name != "b" {
		t.Fatalf("got %+v", repos)
	}
}
