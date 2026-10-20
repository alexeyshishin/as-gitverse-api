package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shishin-omsk/as-gitverse-api/internal/bulk"
	"github.com/shishin-omsk/as-gitverse-api/internal/gitverse"
	"github.com/shishin-omsk/as-gitverse-api/internal/output"
)

func newRepoCmd(flags *globalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Manage repositories (create, list, collaborators, commit-stats), including bulk operations",
	}
	cmd.AddCommand(newRepoCreateCmd(flags))
	cmd.AddCommand(newRepoListCmd(flags))
	cmd.AddCommand(newCollaboratorCmd(flags))
	cmd.AddCommand(newCommitStatsCmd(flags))
	return cmd
}

func newRepoCreateCmd(flags *globalFlags) *cobra.Command {
	var names []string
	var file string
	var private bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create one or more repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			targets, err := resolveNames(names, file)
			if err != nil {
				return err
			}
			client, err := newClient(flags)
			if err != nil {
				return err
			}
			return runBulk(cmd.Context(), flags, targets, func(ctx context.Context, name string) error {
				_, err := client.CreateRepo(ctx, gitverse.CreateRepoRequest{Name: name, Private: private})
				return err
			})
		},
	}
	cmd.Flags().StringSliceVar(&names, "names", nil, "repository names to create (repeatable or comma-separated)")
	cmd.Flags().StringVar(&file, "file", "", "file with one repository name per line")
	cmd.Flags().BoolVar(&private, "private", false, "create as private repositories")
	return cmd
}

func newRepoListCmd(flags *globalFlags) *cobra.Command {
	var org string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List the authenticated user's repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(flags)
			if err != nil {
				return err
			}
			var repos []gitverse.Repository
			if org != "" {
				repos, err = client.ListOrgRepos(cmd.Context(), org)
			} else {
				repos, err = client.ListRepos(cmd.Context())
			}
			if err != nil {
				return err
			}
			rows := make([][]string, len(repos))
			for i, r := range repos {
				rows[i] = []string{r.FullName, visibility(r.Private), r.HTMLURL}
			}
			return newWriter(flags).WriteTable([]string{"NAME", "VISIBILITY", "URL"}, rows)
		},
	}
	cmd.Flags().StringVar(&org, "org", "", "list repositories of this organization instead of the authenticated user")
	return cmd
}

func resolveNames(names []string, file string) ([]string, error) {
	all := append([]string{}, names...)

	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", file, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				all = append(all, line)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("read %s: %w", file, err)
		}
	}

	if len(all) == 0 {
		return nil, fmt.Errorf("no repository names given: use --names or --file")
	}
	return all, nil
}

func runBulk(ctx context.Context, flags *globalFlags, names []string, fn func(context.Context, string) error) error {
	results := bulk.Run(ctx, names, flags.concurrency, fn)

	itemResults := make([]output.ItemResult, len(results))
	failed := 0
	for i, r := range results {
		itemResults[i] = output.ItemResult{Name: r.Item, OK: r.Err == nil}
		if r.Err != nil {
			itemResults[i].Error = r.Err.Error()
			failed++
		}
	}

	if err := newWriter(flags).WriteBulkResults(itemResults); err != nil {
		return err
	}
	if failed > 0 {
		return fmt.Errorf("%d of %d operations failed", failed, len(results))
	}
	return nil
}

func visibility(private bool) string {
	if private {
		return "private"
	}
	return "public"
}
