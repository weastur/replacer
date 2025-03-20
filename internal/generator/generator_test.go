package generator

import (
	"os"
	"regexp"
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
		if err == nil || err.Error() != "failed to get file info: stat nonexistent.txt: no such file or directory" {
			t.Fatalf("expected file not found error, got %v", err)
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
}
