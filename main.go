package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Lint struct {
		Files     []string `arg:"" name:"files" help:"Files or folders to lint"`
		Excluedes []string `short:"e" long:"exclude" help:"Files or folders to exclude"`
	} `cmd:"lint" help:"Lint files or folders"`

	Config string `long:"config" help:"Path to config file"`
}

func lint(ctx context.Context, targets []string, excludes []string, rules map[string]Rule) {
	for _, target := range targets {
		filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
			// TODO: simple string match might not be enough
			for _, exclude := range excludes {
				if path == exclude {
					return nil
				}
			}

			if err != nil {
				log.Fatal("unable to walk through file", err)
			}

			if info.IsDir() {
				return nil
			}

			// TODO: probably we could extract out parser and
			// query creating out to a global one
			for _, rule := range rules {
				query := rule.Query

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

				runLint(ctx, LanguageMap[rule.Language].TSLang, path, query, rule.Message)
			}

			return nil
		})
	}
}

func main() {
	kctx := kong.Parse(&CLI)
	switch kctx.Command() {
	case "lint <files>":
		config, err := parseConfig(CLI.Config)
		if err != nil {
			log.Fatal("unable to parse config", err)
		}

		if len(config.Rules) == 0 {
			log.Fatal("no rules found in config")
		}

		lint(context.Background(), CLI.Lint.Files, CLI.Lint.Excluedes, config.Rules)
	default:
		panic(kctx.Command())
	}
}
