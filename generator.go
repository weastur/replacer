package main

import (
	"fmt"
	"os"
)

func Run(cfg *Config) error {
	filePath := os.Getenv("GOFILE")
	if filePath == "" {
		PrintVerbose("GOFILE environment variable is not set")

		return nil
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	filePerms := fileInfo.Mode().Perm()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	updatedContent := string(content)
	for _, rule := range cfg.Rules {
		updatedContent = rule.Regex.ReplaceAllString(updatedContent, rule.Repl)
	}

	err = os.WriteFile(filePath, []byte(updatedContent), filePerms)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
