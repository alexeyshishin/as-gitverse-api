package gitverse

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
)

const maxPages = 1000

var linkEntryRe = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)

func parseLinkHeader(h http.Header) map[string]string {
	out := map[string]string{}
	for _, part := range linkEntryRe.FindAllStringSubmatch(h.Get("Link"), -1) {
		out[part[2]] = part[1]
	}
	return out
}

// nextPageURL resolves the next page against currentURL's host/path, taking only the
// "page" query parameter from the Link header's next entry. Gitverse's Link header
// points at a different (web, non-API) host, so the URL itself can't be followed as-is.
func nextPageURL(currentURL string, h http.Header) (string, bool) {
	next, ok := parseLinkHeader(h)["next"]
	if !ok {
		return "", false
	}

	nextParsed, err := url.Parse(next)
	if err != nil {
		return "", false
	}
	page := nextParsed.Query().Get("page")
	if page == "" {
		return "", false
	}

	current, err := url.Parse(currentURL)
	if err != nil {
		return "", false
	}
	q := current.Query()
	q.Set("page", page)
	current.RawQuery = q.Encode()
	return current.String(), true
}

func fetchAllPages[T any](ctx context.Context, c *Client, startURL string) ([]T, error) {
	var all []T
	url := startURL

	for page := 0; page < maxPages; page++ {
		var pageItems []T
		headers, err := c.doURL(ctx, "GET", url, nil, &pageItems)
		if err != nil {
			return nil, err
		}
		all = append(all, pageItems...)

		next, ok := nextPageURL(url, headers)
		if !ok {
			break
		}
		url = next
	}

	return all, nil
}
