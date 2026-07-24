package main

import (
	"errors"
	"os"
	"strings"
)

// loadDotEnv reads KEY=VALUE pairs from the file at path and applies them
// via os.Setenv, without overriding any variable already set in the real
// environment (so a real env var always wins over the .env file). It's a
// no-op, not an error, if the file doesn't exist — deployed environments
// (e.g. systemd) are expected to set real env vars directly.
func loadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}
