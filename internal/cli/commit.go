package cli

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func newCommitStatsCmd(flags *globalFlags) *cobra.Command {
	var owner, repo, since, until string

	cmd := &cobra.Command{
		Use:   "commit-stats",
		Short: "Show commit counts per author for a repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			var sinceTime, untilTime time.Time
			if since != "" {
				t, err := time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("invalid --since %q: expected date as YYYY-MM-DD: %w", since, err)
				}
				sinceTime = t
			}
			if until != "" {
				t, err := time.Parse("2006-01-02", until)
				if err != nil {
					return fmt.Errorf("invalid --until %q: expected date as YYYY-MM-DD: %w", until, err)
				}
				// make --until inclusive of the whole day
				untilTime = t.Add(24 * time.Hour)
			}

			client, err := newClient(flags)
			if err != nil {
				return err
			}
			commits, err := client.ListCommits(cmd.Context(), owner, repo)
			if err != nil {
				return err
			}

			counts := map[string]int{}
			for _, c := range commits {
				at := c.AuthoredAt()
				if !sinceTime.IsZero() && at.Before(sinceTime) {
					continue
				}
				if !untilTime.IsZero() && at.After(untilTime) {
					continue
				}
				counts[c.AuthorName()]++
			}

			authors := make([]string, 0, len(counts))
			for a := range counts {
				authors = append(authors, a)
			}
			sort.Slice(authors, func(i, j int) bool { return counts[authors[i]] > counts[authors[j]] })

			rows := make([][]string, len(authors))
			for i, a := range authors {
				rows[i] = []string{a, strconv.Itoa(counts[a])}
			}
			return newWriter(flags).WriteTable([]string{"AUTHOR", "COMMITS"}, rows)
		},
	}
	cmd.Flags().StringVar(&owner, "owner", "", "repository owner")
	cmd.Flags().StringVar(&repo, "repo", "", "repository name")
	cmd.Flags().StringVar(&since, "since", "", "only count commits on or after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "only count commits on or before this date (YYYY-MM-DD)")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")
	return cmd
}
