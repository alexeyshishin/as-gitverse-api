package config

import (
	"os"
	"path/filepath"
	"testing"
)

func testConfigDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CONFIG_HOME", dir)
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("UserConfigDir: %v", err)
	}
	return configDir
}

func TestResolveToken_FlagWins(t *testing.T) {
	testConfigDir(t)
	t.Setenv("GITVERSE_TOKEN", "env-token")
	tok, err := ResolveToken("flag-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "flag-token" {
		t.Fatalf("got %q", tok)
	}
}

func TestResolveToken_EnvWinsOverFile(t *testing.T) {
	configDir := testConfigDir(t)
	writeCredentials(t, configDir, "file-token")

	t.Setenv("GITVERSE_TOKEN", "env-token")
	tok, err := ResolveToken("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "env-token" {
		t.Fatalf("got %q", tok)
	}
}

func TestResolveToken_FallsBackToFile(t *testing.T) {
	configDir := testConfigDir(t)
	t.Setenv("GITVERSE_TOKEN", "")
	writeCredentials(t, configDir, "file-token")

	tok, err := ResolveToken("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "file-token" {
		t.Fatalf("got %q", tok)
	}
}

func TestResolveToken_NoneSet(t *testing.T) {
	testConfigDir(t)
	t.Setenv("GITVERSE_TOKEN", "")

	if _, err := ResolveToken(""); err == nil {
		t.Fatalf("expected error")
	}
}

func TestWriteToken_SetsPermissions(t *testing.T) {
	configDir := testConfigDir(t)

	if err := WriteToken("secret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path := filepath.Join(configDir, credentialsDir, credentialsFile)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("got perm %v", info.Mode().Perm())
	}
}

func writeCredentials(t *testing.T, configDir, token string) {
	t.Helper()
	dir := filepath.Join(configDir, credentialsDir)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, credentialsFile), []byte(token), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}
