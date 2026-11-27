package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	file := filepath.Join(t.TempDir(), ".env")
	content := `
# a comment
TOOLS_TEST_A=hello
TOOLS_TEST_B="quoted value"
TOOLS_TEST_C='single quoted'
TOOLS_TEST_EMPTY=
not a key value line
=missing key
`
	if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TOOLS_TEST_A", "preset")

	if err := Load(file); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := os.Getenv("TOOLS_TEST_A"); got != "preset" {
		t.Errorf("TOOLS_TEST_A = %q, want %q (existing variables must not be overridden)", got, "preset")
	}
	if got := os.Getenv("TOOLS_TEST_B"); got != "quoted value" {
		t.Errorf("TOOLS_TEST_B = %q, want %q", got, "quoted value")
	}
	if got := os.Getenv("TOOLS_TEST_C"); got != "single quoted" {
		t.Errorf("TOOLS_TEST_C = %q, want %q", got, "single quoted")
	}
	if got, ok := os.LookupEnv("TOOLS_TEST_EMPTY"); !ok || got != "" {
		t.Errorf("TOOLS_TEST_EMPTY = %q (set=%v), want an empty set variable", got, ok)
	}
}

func TestLoadMissingFile(t *testing.T) {
	if err := Load(filepath.Join(t.TempDir(), "missing.env")); err == nil {
		t.Error("Load() with a missing file: want an error, got nil")
	}
}

func TestGet(t *testing.T) {
	t.Setenv("TOOLS_TEST_GET", "value")

	if got := Get("TOOLS_TEST_GET", "fallback"); got != "value" {
		t.Errorf("Get() = %q, want %q", got, "value")
	}
	if got := Get("TOOLS_TEST_UNSET", "fallback"); got != "fallback" {
		t.Errorf("Get() = %q, want %q", got, "fallback")
	}
}
