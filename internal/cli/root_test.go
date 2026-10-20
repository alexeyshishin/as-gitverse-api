package cli

import "testing"

func TestNewRootCmd_RegistersSubcommands(t *testing.T) {
	root := NewRootCmd()

	want := map[string]bool{"repo": false}
	for _, c := range root.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("missing subcommand %q", name)
		}
	}
}

func TestRepoCmd_RegistersCreateListCollaboratorsCommitStats(t *testing.T) {
	root := NewRootCmd()
	repoCmd, _, err := root.Find([]string{"repo"})
	if err != nil {
		t.Fatalf("find repo: %v", err)
	}

	want := map[string]bool{"create": false, "list": false, "collaborators": false, "commit-stats": false}
	for _, c := range repoCmd.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("missing repo subcommand %q", name)
		}
	}
}
