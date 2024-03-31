package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	sitter "github.com/smacker/go-tree-sitter"
)

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("unable to open file", err)
	}

	defer f.Close()

	return io.ReadAll(f)
}

func runLint(
	ctx context.Context,
	lang *sitter.Language,
	path string,
	query string,
	message string,
) error {
	sourceCode, err := readFile(path)
	if err != nil {
		return errors.Wrap(err, "unable to read file")
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, sourceCode)
	if err != nil {
		return errors.Wrap(err, "unable to parse file")
	}

	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		return errors.Wrap(err, "unable to create query")
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			// Only use "lint" capture name
			if q.CaptureNameForId(c.Index) != "lint" {
				continue
			}

			msg := strings.Replace(message, "{}", c.Node.Content(sourceCode), -1)
			output := fmt.Sprintf(
				"%s:%d:%d:%d %s",
				path,
				c.Node.StartPoint().Row,
				c.Node.StartPoint().Column,
				c.Node.EndPoint().Column,
				msg,
			)

			// Should we move the print to somewhere else?
			fmt.Println(output)
		}
	}

	return nil
}
