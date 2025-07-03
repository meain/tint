package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Lint struct {
		Files []string `arg:"" name:"files" help:"Files or folders to lint" default:"."`
		Rules []string `short:"r" long:"rule" help:"Rules to check"`
	} `cmd:"lint" help:"Lint files or folders"`

	ValidateConfig struct{} `cmd:"validate-config" help:"Validate config file"`

	Config string `long:"config" help:"Path to config file"`
}

func lint(ctx context.Context, targets []string, rules map[string]Rule, config Config) (int, int) {
	errCount := 0
	fileCount := 0

	for _, target := range targets {
		filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal("unable to walk through file", err)
			}

			// Check if this file should be processed based on config patterns
			if !info.IsDir() && !shouldProcessPath(path, config.Include, config.Exclude) {
				return nil
			}

			fileCount += 1

			// TODO: probably we could extract out parser and
			// query creating out to a global one
			for _, rule := range rules {
				query := rule.Query

				if info.IsDir() {
					continue
				}

				// Check if this file should be processed for this specific rule
				if !shouldProcessPath(path, rule.Include, rule.Exclude) {
					continue
				}

				found := false
				for _, ext := range LanguageMap[rule.Language].Extensions {
					if filepath.Ext(path) == "."+ext {
						found = true
						break
					}
				}

				if !found {
					continue
				}

				// TODO: huge optimization potential
				// We currently parse the file multiple times(once for
				// each rule), no multithreading, no caching
				// We also load the grammars multiple times
				// File is read multiple times from disk
				count, err := runLint(ctx, LanguageMap[rule.Language].TSLang, path, query, rule.Message)
				if err != nil {
					log.Fatal("unable to run lint", err)
				}

				errCount += count
			}

			return nil
		})
	}

	return fileCount, errCount
}

func main() {
	kctx := kong.Parse(&CLI)
	switch kctx.Command() {
	case "lint":
		fallthrough
	case "lint <files>":
		config, err := parseConfig(CLI.Config)
		if err != nil {
			log.Fatal("unable to parse config", err)
		}

		start := time.Now()

		selectedRules := map[string]Rule{}
		if len(CLI.Lint.Rules) == 0 {
			selectedRules = config.Rules
		} else {
			for _, rn := range CLI.Lint.Rules {
				if _, ok := config.Rules[rn]; !ok {
					log.Fatalf("rule %q not found in config", rn)
				}

				selectedRules[rn] = config.Rules[rn]
			}
		}

		fileCount, errCount := lint(context.Background(), CLI.Lint.Files, selectedRules, config)

		// No news is good news
		if errCount == 0 {
			return
		}

		fmt.Fprintf(
			os.Stderr,
			"Found %d issues from %d files using %d rules in %s\n",
			errCount,
			fileCount,
			len(selectedRules),
			time.Since(start).Round(time.Second),
		)

		if errCount > 0 {
			os.Exit(1)
		}
	case "validate-config":
		_, err := parseConfig(CLI.Config)
		if err != nil {
			log.Fatal("unable to parse config: ", err)
		}

		fmt.Println("Config file looks OK")
	default:
		panic(kctx.Command())
	}
}
