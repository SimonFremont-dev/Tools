package command

import (
	"os"
	"strings"
	"testing"
)

func TestOutput(t *testing.T) {
	out, err := Output("go", "version")
	if err != nil {
		t.Fatalf("Output() error = %v", err)
	}
	if !strings.Contains(out, "go") {
		t.Errorf("Output() = %q, want it to contain %q", out, "go")
	}
}

func TestOutputError(t *testing.T) {
	if _, err := Output("tools-definitely-not-a-real-binary"); err == nil {
		t.Error("Output() with a missing binary: want an error, got nil")
	}
}

func TestRun(t *testing.T) {
	if err := Run("go", "version"); err != nil {
		t.Errorf("Run() error = %v", err)
	}
}

func TestRunIn(t *testing.T) {
	dir := t.TempDir()
	if err := RunIn(dir, "go", "version"); err != nil {
		t.Errorf("RunIn() error = %v", err)
	}
}

func TestOutputInEnv(t *testing.T) {
	cmd := newCmd("", []string{"TOOLS_TEST_CMD=hello"}, "go", "version")
	found := false
	for _, e := range cmd.Env {
		if e == "TOOLS_TEST_CMD=hello" {
			found = true
		}
	}
	if !found {
		t.Error("newCmd() did not append the extra environment variable")
	}
	if cmd.Dir != "" {
		t.Errorf("newCmd().Dir = %q, want empty", cmd.Dir)
	}
}

func TestNewCmdDir(t *testing.T) {
	cmd := newCmd("/tmp", nil, "go")
	if cmd.Dir != "/tmp" {
		t.Errorf("newCmd().Dir = %q, want %q", cmd.Dir, "/tmp")
	}
	if cmd.Env != nil {
		t.Errorf("newCmd().Env = %v, want nil (inherit)", cmd.Env)
	}
	if _, ok := os.LookupEnv("TOOLS_TEST_UNUSED"); ok {
		t.Error("unexpected TOOLS_TEST_UNUSED in environment")
	}
}
