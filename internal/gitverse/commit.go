package gitverse

import (
	"context"
	"time"
)

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name string    `json:"name"`
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
}

func (cm Commit) AuthorName() string {
	if cm.Author.Login != "" {
		return cm.Author.Login
	}
	return cm.Commit.Author.Name
}

func (cm Commit) AuthoredAt() time.Time {
	return cm.Commit.Author.Date
}

// ListCommits fetches the full commit history unfiltered. The server-side
// since/until query params are not used: Gitverse's API silently truncates
// pagination when they're set (observed live — a repo with 490 commits
// spanning ~7 weeks returns only ~100 commits when `until` is applied,
// regardless of how far back `until` reaches). Callers filter by date
// client-side instead.
func (c *Client) ListCommits(ctx context.Context, owner, repo string) ([]Commit, error) {
	return fetchAllPages[Commit](ctx, c, c.baseURL+buildPath("repos", owner, repo, "commits"))
}
