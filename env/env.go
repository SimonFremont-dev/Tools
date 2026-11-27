// Package env loads KEY=VALUE pairs from env files into the process environment.
package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load reads KEY=VALUE pairs from the given files and sets them as
// environment variables. Variables already present in the environment are
// not overridden. When no file is given, ".env" is used.
func Load(filenames ...string) error {
	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, name := range filenames {
		err := loadFile(name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get returns the value of the environment variable named by key, or
// fallback when the variable is not set.
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func loadFile(name string) error {
	file, err := os.Open(name) // #nosec G304 -- the caller chooses the env file path
	if err != nil {
		return fmt.Errorf("env: open %s: %w", name, err)
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++

		key, value, ok := parseLine(scanner.Text())
		if !ok {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			err = os.Setenv(key, value)
			if err != nil {
				return fmt.Errorf("env: set %s (line %d of %s): %w", key, line, name, err)
			}
		}
	}

	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("env: read %s: %w", name, err)
	}

	return nil
}

func parseLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)

	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}

	key, value, found := strings.Cut(line, "=")
	if !found {
		return "", "", false
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return "", "", false
	}

	value = strings.Trim(strings.TrimSpace(value), `"'`)

	return key, value, true
}
