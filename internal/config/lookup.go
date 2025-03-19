package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotFound = errors.New("config file not found")

func Lookup(path string) (string, error) {
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		fmt.Printf("Config file not found at %s\n", path)

		return "", ErrNotFound
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %s\n", err)

		return "", ErrNotFound
	}

	for {
		ymlPath := filepath.Join(currentDir, ".replacer.yml")
		yamlPath := filepath.Join(currentDir, ".replacer.yaml")
		goModPath := filepath.Join(currentDir, "go.mod")

		if _, err := os.Stat(ymlPath); err == nil {
			return ymlPath, nil
		}

		if _, err := os.Stat(yamlPath); err == nil {
			return yamlPath, nil
		}

		if _, err := os.Stat(goModPath); err == nil {
			fmt.Println("Reached directory containing go.mod, no config file found")

			return "", ErrNotFound
		}

		parentDir := filepath.Dir(currentDir)

		if parentDir == currentDir {
			fmt.Println("Reached root directory, no config file found")

			break
		}

		currentDir = parentDir
	}

	return "", ErrNotFound
}
