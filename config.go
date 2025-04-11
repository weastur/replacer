package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/goccy/go-yaml"
)

var ErrNotFound = errors.New("config file not found")

type MissingRequiredFieldError struct {
	Field string
}

func (e *MissingRequiredFieldError) Error() string {
	return "missing required field " + e.Field
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Regex *regexp.Regexp `yaml:"regex"`
	Repl  string         `yaml:"repl"`
}

func (p *Rule) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw struct {
		Regex *string `yaml:"regex"`
		Repl  *string `yaml:"repl"`
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.Regex == nil {
		return &MissingRequiredFieldError{"regex"}
	}

	if raw.Repl == nil {
		return &MissingRequiredFieldError{"repl"}
	}

	compiledRegex, err := regexp.Compile(*raw.Regex)
	if err != nil {
		return fmt.Errorf("error compiling regex: %w", err)
	}

	p.Regex = compiledRegex
	p.Repl = *raw.Repl

	return nil
}

func LookupConfig(path string) (string, error) { //nolint:cyclop
	if path != "" {
		var err error

		if _, err = os.Stat(path); err == nil {
			return path, nil
		}

		return "", fmt.Errorf("config file not found at %s: %w", path, err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	for {
		ymlPath := filepath.Join(currentDir, ".replacer.yml")
		yamlPath := filepath.Join(currentDir, ".replacer.yaml")
		goModPath := filepath.Join(currentDir, "go.mod")

		if _, err := os.Stat(ymlPath); err == nil {
			return ymlPath, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("error checking file %s: %w", ymlPath, err)
		}

		if _, err := os.Stat(yamlPath); err == nil {
			return yamlPath, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("error checking file %s: %w", yamlPath, err)
		}

		if _, err := os.Stat(goModPath); err == nil {
			PrintVerbose("Reached directory containing go.mod, no config file found")

			return "", ErrNotFound
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("error checking file %s: %w", goModPath, err)
		}

		parentDir := filepath.Dir(currentDir)

		if parentDir == currentDir {
			PrintVerbose("Reached root directory, no config file found")

			break
		}

		currentDir = parentDir
	}

	return "", ErrNotFound
}

func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
