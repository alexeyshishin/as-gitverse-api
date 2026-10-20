package cli

import "testing"

func TestCollaboratorsCmd_RegistersListAndAdd(t *testing.T) {
	root := NewRootCmd()
	repoCmd, _, err := root.Find([]string{"repo"})
	if err != nil {
		t.Fatalf("find repo: %v", err)
	}
	collabCmd, _, err := repoCmd.Find([]string{"collaborators"})
	if err != nil {
		t.Fatalf("find collaborators: %v", err)
	}

	want := map[string]bool{"list": false, "add": false}
	for _, c := range collabCmd.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Fatalf("missing collaborators subcommand %q", name)
		}
	}
}

func TestCollaboratorsListCmd_RequiresOwnerAndRepo(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"repo", "collaborators", "list"})
	if err := root.Execute(); err == nil {
		t.Fatalf("expected error for missing --owner/--repo")
	}
}

func TestCollaboratorsAddCmd_RequiresPermission(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"repo", "collaborators", "add", "--owner", "o", "--repo", "r", "--names", "u"})
	if err := root.Execute(); err == nil {
		t.Fatalf("expected error for missing --permission")
	}
}

func TestCommitStatsCmd_RequiresOwnerAndRepo(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"repo", "commit-stats"})
	if err := root.Execute(); err == nil {
		t.Fatalf("expected error for missing --owner/--repo")
	}
}

func TestRepoListCmd_HasOrgFlag(t *testing.T) {
	root := NewRootCmd()
	repoCmd, _, err := root.Find([]string{"repo"})
	if err != nil {
		t.Fatalf("find repo: %v", err)
	}
	listCmd, _, err := repoCmd.Find([]string{"list"})
	if err != nil {
		t.Fatalf("find list: %v", err)
	}
	if listCmd.Flags().Lookup("org") == nil {
		t.Fatalf("expected --org flag on repo list")
	}
}
