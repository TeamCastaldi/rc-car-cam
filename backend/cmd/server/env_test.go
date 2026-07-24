package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnv_SetsSimpleKeyValuePairs(t *testing.T) {
	_ = os.Unsetenv("TEST_LOAD_DOTENV_A")
	_ = os.Unsetenv("TEST_LOAD_DOTENV_B")
	t.Cleanup(func() {
		_ = os.Unsetenv("TEST_LOAD_DOTENV_A")
		_ = os.Unsetenv("TEST_LOAD_DOTENV_B")
	})

	path := writeTempEnvFile(t, "TEST_LOAD_DOTENV_A=one\nTEST_LOAD_DOTENV_B=two\n")

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv: %v", err)
	}

	if got := os.Getenv("TEST_LOAD_DOTENV_A"); got != "one" {
		t.Errorf("TEST_LOAD_DOTENV_A = %q, want %q", got, "one")
	}
	if got := os.Getenv("TEST_LOAD_DOTENV_B"); got != "two" {
		t.Errorf("TEST_LOAD_DOTENV_B = %q, want %q", got, "two")
	}
}

func TestLoadDotEnv_SkipsBlankLinesAndComments(t *testing.T) {
	_ = os.Unsetenv("TEST_LOAD_DOTENV_C")
	t.Cleanup(func() { _ = os.Unsetenv("TEST_LOAD_DOTENV_C") })

	path := writeTempEnvFile(t, "# a comment\n\nTEST_LOAD_DOTENV_C=three\n")

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv: %v", err)
	}

	if got := os.Getenv("TEST_LOAD_DOTENV_C"); got != "three" {
		t.Errorf("TEST_LOAD_DOTENV_C = %q, want %q", got, "three")
	}
}

func TestLoadDotEnv_DoesNotOverrideRealEnv(t *testing.T) {
	t.Setenv("TEST_LOAD_DOTENV_D", "real-value")

	path := writeTempEnvFile(t, "TEST_LOAD_DOTENV_D=from-file\n")

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv: %v", err)
	}

	if got := os.Getenv("TEST_LOAD_DOTENV_D"); got != "real-value" {
		t.Errorf("TEST_LOAD_DOTENV_D = %q, want %q (real env should win)", got, "real-value")
	}
}

func TestLoadDotEnv_MissingFileIsNotAnError(t *testing.T) {
	if err := loadDotEnv(filepath.Join(t.TempDir(), "does-not-exist.env")); err != nil {
		t.Errorf("expected no error for a missing .env file, got: %v", err)
	}
}

func TestLoadDotEnv_SkipsMalformedLines(t *testing.T) {
	_ = os.Unsetenv("TEST_LOAD_DOTENV_E")
	t.Cleanup(func() { _ = os.Unsetenv("TEST_LOAD_DOTENV_E") })

	path := writeTempEnvFile(t, "this line has no equals sign\nTEST_LOAD_DOTENV_E=five\n")

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("loadDotEnv: %v", err)
	}

	if got := os.Getenv("TEST_LOAD_DOTENV_E"); got != "five" {
		t.Errorf("TEST_LOAD_DOTENV_E = %q, want %q", got, "five")
	}
}

func writeTempEnvFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write temp .env file: %v", err)
	}
	return path
}
