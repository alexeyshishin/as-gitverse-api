package gitverse

import "context"

type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
}

type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Private     bool   `json:"private"`
}

func (c *Client) CreateRepo(ctx context.Context, req CreateRepoRequest) (*Repository, error) {
	var repo Repository
	if _, err := c.do(ctx, "POST", "/user/repos", req, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (c *Client) ListRepos(ctx context.Context) ([]Repository, error) {
	return fetchAllPages[Repository](ctx, c, c.baseURL+"/user/repos")
}

func (c *Client) ListOrgRepos(ctx context.Context, org string) ([]Repository, error) {
	return fetchAllPages[Repository](ctx, c, c.baseURL+buildPath("orgs", org, "repos"))
}
