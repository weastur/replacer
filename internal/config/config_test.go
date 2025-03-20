package config

import (
	"errors"
	"os"
	"testing"

	"github.com/goccy/go-yaml"
)

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

		if _, err = Load(tempFile.Name()); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("File does not exist", func(t *testing.T) {
		_, err := Load("nonexistent.yml")
		if err == nil {
			t.Fatal("Expected error, got nil")
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

		_, err = Load(tempFile.Name())
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}
