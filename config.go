package main

import (
	"fmt"
	"os"

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
}

type Config struct {
	Rules map[string]Rule
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
