package gitverse

import (
	"context"
	"fmt"
)

type Collaborator struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Permissions struct {
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
		Admin bool `json:"admin"`
	} `json:"permissions"`
}

func (c Collaborator) PermissionLabel() string {
	switch {
	case c.Permissions.Admin:
		return "admin"
	case c.Permissions.Push:
		return "write"
	case c.Permissions.Pull:
		return "read"
	default:
		return ""
	}
}

type addCollaboratorRequest struct {
	Permission string `json:"permission"`
}

var ValidCollaboratorPermissions = []string{"read", "write", "admin"}

func ValidateCollaboratorPermission(permission string) error {
	for _, p := range ValidCollaboratorPermissions {
		if p == permission {
			return nil
		}
	}
	return fmt.Errorf("invalid permission %q, must be one of %v", permission, ValidCollaboratorPermissions)
}

func (c *Client) ListCollaborators(ctx context.Context, owner, repo string) ([]Collaborator, error) {
	return fetchAllPages[Collaborator](ctx, c, c.baseURL+buildPath("repos", owner, repo, "collaborators"))
}

func (c *Client) AddCollaborator(ctx context.Context, owner, repo, username, permission string) error {
	if err := ValidateCollaboratorPermission(permission); err != nil {
		return err
	}
	_, err := c.do(ctx, "PUT", buildPath("repos", owner, repo, "collaborators", username), addCollaboratorRequest{Permission: permission}, nil)
	return err
}
