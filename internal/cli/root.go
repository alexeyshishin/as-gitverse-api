package cli

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/shishin-omsk/as-gitverse-api/internal/config"
	"github.com/shishin-omsk/as-gitverse-api/internal/gitverse"
	"github.com/shishin-omsk/as-gitverse-api/internal/output"
)

var cmdOut io.Writer = os.Stdout

type globalFlags struct {
	token       string
	jsonOutput  bool
	concurrency int
	verbose     bool
}

func NewRootCmd() *cobra.Command {
	flags := &globalFlags{}

	root := &cobra.Command{
		Use:           "gitverse",
		Short:         "Bulk management of Gitverse repositories",
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	root.PersistentFlags().StringVar(&flags.token, "token", "", "Gitverse Personal Access Token (overrides GITVERSE_TOKEN and the credentials file)")
	root.PersistentFlags().BoolVar(&flags.jsonOutput, "json", false, "output as JSON instead of a table")
	root.PersistentFlags().IntVar(&flags.concurrency, "concurrency", 5, "max concurrent requests for bulk operations")
	root.PersistentFlags().BoolVar(&flags.verbose, "verbose", false, "verbose logging (never includes the token)")

	root.AddCommand(newRepoCmd(flags))

	return root
}

func newClient(flags *globalFlags) (*gitverse.Client, error) {
	token, err := config.ResolveToken(flags.token)
	if err != nil {
		return nil, err
	}
	client := gitverse.NewClient(token)
	client.SetVerbose(flags.verbose)
	return client, nil
}

func newWriter(flags *globalFlags) output.Writer {
	return output.New(cmdOut, flags.jsonOutput)
}
