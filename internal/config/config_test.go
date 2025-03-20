package config

import (
	"errors"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestRule_UnmarshalYAML(t *testing.T) {
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
		if !errors.As(err, &missingFieldErr) || err.Error() != "missing required field regex" {
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
		if !errors.As(err, &missingFieldErr) || err.Error() != "missing required field repl" {
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

		if err.Error() != "error compiling regex: error parsing regexp: missing closing ]: `[invalid`" {
			t.Errorf("Expected regex compilation error, got %v", err)
		}
	})
}
