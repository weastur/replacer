package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/goccy/go-yaml"
)

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

func Load(filepath string) (*Config, error) {
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
