package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config []CommandDefinition

type CommandDefinition struct {
	Cmd             string `yaml:"Cmd"`
	OpenTag         string `yaml:"OpenTag"`
	CloseTag        string `yaml:"CloseTag"`
	ReplaceOpenTag  string `yaml:"ReplaceOpenTag"`
	ReplaceCloseTag string `yaml:"ReplaceCloseTag"`
	InputPrefix     string `yaml:"InputPrefix"`
	InputPostfix    string `yaml:"InputPostfix"`
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	body, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("cannot open file: %w", err)
	}

	err = yaml.Unmarshal(body, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("cannot decode yaml: %w", err)
	}

	for i := range cfg {
		err := cfg[i].validateAndFix()
		if err != nil {
			return cfg, fmt.Errorf("invalid command definition: %w", err)
		}
	}

	if len(cfg) == 0 {
		return cfg, fmt.Errorf("no command defined")
	}
	return cfg, nil
}

func (c *CommandDefinition) validateAndFix() error {
	if strings.TrimSpace(c.Cmd) == "" { // No Cmd
		return errors.New("no command defined")
	}
	if strings.TrimSpace(c.OpenTag) == "" { // No OpenTag
		return errors.New("no open tag defined")
	}

	if c.CloseTag == "" {
		c.CloseTag = c.OpenTag
	}

	return nil
}

func (c *CommandDefinition) DetectCloseCommand(line []byte) (mayCloseCommand bool) {
	if c == nil {
		return false
	}

	return strings.HasPrefix(string(line), c.CloseTag)
}

func (c Config) DetectOpenCommand(line []byte) (cmd *CommandDefinition, detected bool) {
	for _, cmd := range c {
		if strings.HasPrefix(string(line), cmd.OpenTag) {
			return &cmd, true
		}
	}

	return nil, false
}
