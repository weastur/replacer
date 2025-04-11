package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestLookup(t *testing.T) {
	t.Run("Valid file path", func(t *testing.T) {
		tempFile, err := os.CreateTemp(t.TempDir(), "test-config.yml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		path, err := LookupConfig(tempFile.Name())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if path != tempFile.Name() {
			t.Errorf("Expected path %s, got %s", tempFile.Name(), path)
		}
	})

	t.Run("Invalid file path", func(t *testing.T) {
		_, err := LookupConfig("/invalid/path/to/config.yml")
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Error in getwd", func(t *testing.T) {
		tempDir := t.TempDir()

		tempFile := filepath.Join(tempDir, ".replacer.yml")
		if err := os.WriteFile(tempFile, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		originalDir, _ := os.Getwd()
		defer t.Chdir(originalDir)
		t.Chdir(tempDir)

		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Failed to remove directory: %v", err)
		}

		_, err := LookupConfig("")

		if err == nil {
			t.Error("Expected error")
		}
	})

	for _, fileName := range []string{".replacer.yml", ".replacer.yaml"} {
		t.Run("Search for config file in current directory: "+fileName, func(t *testing.T) {
			tempDir := t.TempDir()

			tempFile := filepath.Join(tempDir, fileName)
			if err := os.WriteFile(tempFile, []byte("test"), 0o644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			originalDir, _ := os.Getwd()
			defer t.Chdir(originalDir)
			t.Chdir(tempDir)

			path, err := LookupConfig("")
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if path != tempFile {
				t.Errorf("Expected path %s, got %s", tempFile, path)
			}
		})
	}

	t.Run("Reach root directory without finding config file", func(t *testing.T) {
		tempDir := t.TempDir()

		originalDir, _ := os.Getwd()
		defer t.Chdir(originalDir)
		t.Chdir(tempDir)

		_, err := LookupConfig("")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	t.Run("Encounter go.mod without finding config file", func(t *testing.T) {
		tempDir := t.TempDir()

		goModFile := filepath.Join(tempDir, "go.mod")
		if err := os.WriteFile(goModFile, []byte("module test"), 0o644); err != nil {
			t.Fatalf("Failed to create go.mod file: %v", err)
		}

		originalDir, _ := os.Getwd()
		defer t.Chdir(originalDir)
		t.Chdir(tempDir)

		_, err := LookupConfig("")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}

func TestRule_UnmarshalYAML(t *testing.T) {
	t.Run("Invalid yaml", func(t *testing.T) {
		data := "a => b"

		var rule Rule

		err := yaml.Unmarshal([]byte(data), &rule)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if errors.Is(err, &yaml.UnexpectedNodeTypeError{}) {
			t.Errorf("Expected yaml parsing error, got %v", err)
		}
	})

	t.Run("Valid regex and repl", func(t *testing.T) {
		data := `
regex: "^test.*"
repl: "replacement"
`

		var rule Rule

		if err := yaml.Unmarshal([]byte(data), &rule); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if rule.Regex == nil || rule.Regex.String() != "^test.*" {
			t.Errorf("Expected regex to be '^test.*', got %v", rule.Regex)
		}

		if rule.Repl != "replacement" {
			t.Errorf("Expected repl to be 'replacement', got %s", rule.Repl)
		}
	})

	t.Run("Missing regex field", func(t *testing.T) {
		data := `
repl: "replacement"
`

		var rule Rule

		err := yaml.Unmarshal([]byte(data), &rule)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		var missingFieldErr *MissingRequiredFieldError
		if !errors.As(err, &missingFieldErr) || missingFieldErr.Field != "regex" {
			t.Errorf("Expected MissingRequiredFieldError for 'regex', got %v", err)
		}

		if !strings.Contains(missingFieldErr.Error(), "regex") {
			t.Errorf("Expected error message to contain 'regex', got %v", missingFieldErr.Error())
		}
	})

	t.Run("Missing repl field", func(t *testing.T) {
		data := `
regex: "^test.*"
`

		var rule Rule

		err := yaml.Unmarshal([]byte(data), &rule)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		var missingFieldErr *MissingRequiredFieldError
		if !errors.As(err, &missingFieldErr) || missingFieldErr.Field != "repl" {
			t.Errorf("Expected MissingRequiredFieldError for 'repl', got %v", err)
		}

		if !strings.Contains(missingFieldErr.Error(), "repl") {
			t.Errorf("Expected error message to contain 'repl', got %v", missingFieldErr.Error())
		}
	})

	t.Run("Invalid regex pattern", func(t *testing.T) {
		data := `
regex: "[invalid"
repl: "replacement"
`

		var rule Rule

		err := yaml.Unmarshal([]byte(data), &rule)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("Valid config file", func(t *testing.T) {
		data := `
rules:
  - regex: "^test.*"
    repl: "replacement"
  - regex: "foo"
    repl: "bar"
`

		tempFile, err := os.CreateTemp(t.TempDir(), "config.yml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		defer os.Remove(tempFile.Name())

		if _, err := tempFile.WriteString(data); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		tempFile.Close()

		if _, err = LoadConfig(tempFile.Name()); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, err := LoadConfig("nonexistent.yml")
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("Expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("Invalid YAML format", func(t *testing.T) {
		data := `
invalid_yaml
`

		tempFile, err := os.CreateTemp(t.TempDir(), "config.yml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}

		defer os.Remove(tempFile.Name())

		if _, err := tempFile.WriteString(data); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}

		tempFile.Close()

		_, err = LoadConfig(tempFile.Name())
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
