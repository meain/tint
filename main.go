package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	sitter "github.com/smacker/go-tree-sitter"
)

var CLI struct {
	Lint struct {
		Files     []string `arg:"" name:"files" help:"Files or folders to lint"`
		Excluedes []string `short:"e" long:"exclude" help:"Files or folders to exclude"`
	} `cmd:"lint" help:"Lint files or folders"`

	Config string `short:"c" long:"config" help:"Path to config file"`
}

// TODO: reorder args
func runRule(ctx context.Context, path string, query string, lang *sitter.Language, message string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("unable to open file", err)
	}

	defer f.Close()

	sourceCode, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("unable to read file", err)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, sourceCode)
	if err != nil {
		log.Fatal("unable to parse", err)
	}

	// Define the query
	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		log.Fatal("unable to create query", err)
	}

	// Execute the query
	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	// Iterate over the matches
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			// Only use "lint" capture name
			if q.CaptureNameForId(c.Index) != "lint" {
				continue
			}

			msg := strings.Replace(message, "{}", c.Node.Content(sourceCode), -1)
			output := "%s:%d:%d:%d %s"
			fmt.Println(fmt.Sprintf(output, path, c.Node.StartPoint().Row, c.Node.StartPoint().Column, c.Node.EndPoint().Column, msg))
		}
	}
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

			// TODO(meain): probably we could extract out parser and
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

				runRule(ctx, path, query, LanguageMap[rule.Language].TSLang, rule.Message)
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

		lint(context.Background(), CLI.Lint.Files, CLI.Lint.Excluedes, config.Rules)
	default:
		panic(kctx.Command())
	}
}
