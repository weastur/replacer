package generator

import (
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/weastur/replacer/internal/config"
)

func TestRun(t *testing.T) {
	t.Run("GOFILE not set", func(t *testing.T) {
		os.Unsetenv("GOFILE")

		err := Run(&config.Config{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Setenv("GOFILE", "nonexistent.txt")

		err := Run(&config.Config{})
		if err == nil || !strings.Contains(err.Error(), "no such file") {
			t.Fatalf("expected file not found error, got %v", err)
		}
	})

	t.Run("File exists but unreadable", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("permission bits not reliable on Windows")
		}

		tmpFile, err := os.CreateTemp(t.TempDir(), "testfile*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Make file unreadable
		if err := os.Chmod(tmpFile.Name(), 0); err != nil {
			t.Fatalf("failed to chmod file: %v", err)
		}

		t.Setenv("GOFILE", tmpFile.Name())

		err = Run(&config.Config{})
		if err == nil || !strings.Contains(err.Error(), "failed to read") {
			t.Fatalf("expected error, got %v", err)
		}
	})

	t.Run("File exists and rules applied", func(t *testing.T) {
		tmpFile, err := os.CreateTemp(t.TempDir(), "testfile*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		initialContent := "hello world"
		if _, err := tmpFile.WriteString(initialContent); err != nil {
			t.Fatalf("failed to write to temp file: %v", err)
		}

		tmpFile.Close()

		t.Setenv("GOFILE", tmpFile.Name())

		mockConfig := &config.Config{
			Rules: []config.Rule{
				{
					Regex: regexp.MustCompile("hello"),
					Repl:  "hi",
				},
				{
					Regex: regexp.MustCompile("world"),
					Repl:  "universe",
				},
			},
		}

		err = Run(mockConfig)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updatedContent, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("failed to read temp file: %v", err)
		}

		expectedContent := "hi universe"
		if string(updatedContent) != expectedContent {
			t.Fatalf("expected content %q, got %q", expectedContent, string(updatedContent))
		}
	})

	t.Run("File exists, rules applied, but can't write back", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("permission bits not reliable on Windows")
		}

		tmpFile, err := os.CreateTemp(t.TempDir(), "testfile*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		t.Setenv("GOFILE", tmpFile.Name())

		// Make file unwritable
		if err := os.Chmod(tmpFile.Name(), 0o400); err != nil {
			t.Fatalf("failed to chmod file: %v", err)
		}

		err = Run(&config.Config{})
		if err == nil || !strings.Contains(err.Error(), "failed to write") {
			t.Fatalf("expected error, got %v", err)
		}
	})
}
