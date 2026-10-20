package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/shishin-omsk/as-gitverse-api/internal/gitverse"
)

func newCollaboratorCmd(flags *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collaborators",
		Short: "Manage repository collaborators",
	}
	cmd.AddCommand(newCollaboratorListCmd(flags))
	cmd.AddCommand(newCollaboratorAddCmd(flags))
	return cmd
}

func newCollaboratorListCmd(flags *globalFlags) *cobra.Command {
	var owner, repo string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repository collaborators",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(flags)
			if err != nil {
				return err
			}
			collaborators, err := client.ListCollaborators(cmd.Context(), owner, repo)
			if err != nil {
				return err
			}
			rows := make([][]string, len(collaborators))
			for i, c := range collaborators {
				rows[i] = []string{c.Login, c.PermissionLabel()}
			}
			return newWriter(flags).WriteTable([]string{"LOGIN", "PERMISSION"}, rows)
		},
	}
	cmd.Flags().StringVar(&owner, "owner", "", "repository owner")
	cmd.Flags().StringVar(&repo, "repo", "", "repository name")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")
	return cmd
}

func newCollaboratorAddCmd(flags *globalFlags) *cobra.Command {
	var names []string
	var file string
	var owner, repo, permission string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add one or more collaborators to a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			targets, err := resolveNames(names, file)
			if err != nil {
				return err
			}
			if err := gitverse.ValidateCollaboratorPermission(permission); err != nil {
				return err
			}
			client, err := newClient(flags)
			if err != nil {
				return err
			}
			return runBulk(cmd.Context(), flags, targets, func(ctx context.Context, name string) error {
				return client.AddCollaborator(ctx, owner, repo, name, permission)
			})
		},
	}
	cmd.Flags().StringSliceVar(&names, "names", nil, "usernames to add as collaborators (repeatable or comma-separated)")
	cmd.Flags().StringVar(&file, "file", "", "file with one username per line")
	cmd.Flags().StringVar(&owner, "owner", "", "repository owner")
	cmd.Flags().StringVar(&repo, "repo", "", "repository name")
	cmd.Flags().StringVar(&permission, "permission", "", "permission to grant: read, write or admin")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")
	_ = cmd.MarkFlagRequired("permission")
	return cmd
}
