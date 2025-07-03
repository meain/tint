package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	sitter "github.com/smacker/go-tree-sitter"
)

type Rule struct {
	// Message is the message to be displayed when the rule is violated
	// {} is replaced with the primary object if present
	Message string

	// Language is the language of the query
	Language Language

	// Query is the tree-sitter query to be run
	// Query can contain multiple queries
	Query string

	// Include patterns specific to this rule
	Include []string `toml:"include,omitempty"`

	// Exclude patterns specific to this rule
	Exclude []string `toml:"exclude,omitempty"`
}

type Config struct {
	Rules   map[string]Rule
	Include []string `toml:"include,omitempty"`
	Exclude []string `toml:"exclude,omitempty"`
}

func validateConfig(config Config) error {
	if len(config.Rules) == 0 {
		return errors.New("no rules found in the config")
	}

	for name, rule := range config.Rules {
		if rule.Message == "" {
			return errors.New("message not found for rule " + name)
		}

		if rule.Language == "" {
			return errors.New("language not found for rule " + name)
		}

		if rule.Query == "" {
			return errors.New("query not found for rule " + name)
		}

		q, err := sitter.NewQuery([]byte(rule.Query), LanguageMap[rule.Language].TSLang)
		if err != nil {
			return errors.Wrap(err, "unable to create query")
		}

		// check if a capture names `region` is available
		found := false
		for i := uint32(0); i < q.CaptureCount(); i++ {
			if q.CaptureNameForId(i) == "region" {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("query for rule %s does not contain a capture named `region`", name)
		}
	}

	return nil
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

	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// shouldProcessPath determines if a path should be processed based on include/exclude patterns
func shouldProcessPath(path string, includePatterns, excludePatterns []string) bool {
	// Convert to clean path for consistent matching
	cleanPath := filepath.Clean(path)

	// If include patterns are specified, the path must match at least one
	if len(includePatterns) > 0 {
		included := false
		for _, pattern := range includePatterns {
			match, err := doublestar.Match(pattern, cleanPath)
			if err != nil {
				// If pattern is invalid, skip it and continue
				continue
			}
			if match {
				included = true
				break
			}
		}
		// If include patterns exist but none matched, exclude the path
		if !included {
			return false
		}
	}

	// Check exclude patterns - if any match, exclude the path
	for _, pattern := range excludePatterns {
		match, err := doublestar.Match(pattern, cleanPath)
		if err != nil {
			// If pattern is invalid, skip it and continue
			continue
		}
		if match {
			return false
		}
	}

	return true
}
