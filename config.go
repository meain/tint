package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Rule struct {
	// Kind is the kind of the rule (eg: error, warning)
	Kind string

	// Message is the message to be displayed when the rule is violated
	// {} is replaced with the primary object if present
	Message string

	// Language is the language of the query
	Language Language

	// Query is the tree-sitter query to be run
	Query string
}

type Config struct {
	Rules map[string]Rule
}

// parseConfig parses the configuration file and returns a Config struct
// Config file should be looked up in .tint.toml if not provided
func parseConfig(config string) (Config, error) {
	if config == "" {
		config = ".tint.toml"
	}

	configData, err := os.ReadFile(config)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := toml.Unmarshal(configData, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
