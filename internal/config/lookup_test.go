package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLookup(t *testing.T) {
	t.Run("Valid file path", func(t *testing.T) {
		tempFile, err := os.CreateTemp(t.TempDir(), "test-config.yml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		path, err := Lookup(tempFile.Name())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if path != tempFile.Name() {
			t.Errorf("Expected path %s, got %s", tempFile.Name(), path)
		}
	})

	t.Run("Invalid file path", func(t *testing.T) {
		_, err := Lookup("/invalid/path/to/config.yml")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got %v", err)
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

			path, err := Lookup("")
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

		_, err := Lookup("")
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

		_, err := Lookup("")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})
}
