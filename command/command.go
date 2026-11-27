package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run executes name with args, wiring the command to the current process
// standard streams. It returns an error if the command cannot be started or
// exits with a non-zero status.
func Run(name string, args ...string) error {
	cmd := newCmd("", nil, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return run(cmd)
}

// RunIn is like Run but executes the command in dir.
func RunIn(dir, name string, args ...string) error {
	cmd := newCmd(dir, nil, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return run(cmd)
}

// Output executes name with args and returns its trimmed combined output.
func Output(name string, args ...string) (string, error) {
	return OutputIn("", nil, name, args...)
}

// OutputIn is like Output but executes the command in dir, with additional
// KEY=VALUE environment variables layered on top of the current ones.
func OutputIn(dir string, env []string, name string, args ...string) (string, error) {
	var buf bytes.Buffer
	cmd := newCmd(dir, env, name, args...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command: %s: %w (output: %s)", cmd.String(), err, strings.TrimSpace(buf.String()))
	}
	return strings.TrimSpace(buf.String()), nil
}

func newCmd(dir string, env []string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	return cmd
}

func run(cmd *exec.Cmd) error {
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command: %s: %w", cmd.String(), err)
	}
	return nil
}
