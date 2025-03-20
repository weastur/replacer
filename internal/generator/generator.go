package generator

import (
	"fmt"
	"os"

	"github.com/weastur/replacer/internal/config"
)

func Run(cfg *config.Config) error {
	filePath := os.Getenv("GOFILE")
	if filePath == "" {
		fmt.Println("GOFILE environment variable is not set")

		return nil
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Failed to get file info: %v\n", err)

		return fmt.Errorf("failed to get file info: %w", err)
	}

	filePerms := fileInfo.Mode().Perm()

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)

		return fmt.Errorf("failed to read file: %w", err)
	}

	updatedContent := string(content)
	for _, rule := range cfg.Rules {
		updatedContent = rule.Regex.ReplaceAllString(updatedContent, rule.Repl)
	}

	err = os.WriteFile(filePath, []byte(updatedContent), filePerms)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)

		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
