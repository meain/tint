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

// TODO: Some way to ignore certain instances for a lint error
// Other option would be to have the user itself form the query so as
// to look for a comment
func runLint(
	ctx context.Context,
	lang *sitter.Language,
	path string,
	query string,
	message string,
) (int, error) {
	errCount := 0

	sourceCode, err := readFile(path)
	if err != nil {
		return 0, errors.Wrap(err, "unable to read file")
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, sourceCode)
	if err != nil {
		return 0, errors.Wrap(err, "unable to parse file")
	}

	q, err := sitter.NewQuery([]byte(query), lang)
	if err != nil {
		return 0, errors.Wrap(err, "unable to create query")
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
			if q.CaptureNameForId(c.Index) != "region" {
				continue
			}

			// TODO: user should be able to use any capture here and a
			// separate one for the region to mark
			// NOTE: Start and end could be on separate lines
			msg := strings.Replace(message, "{}", c.Node.Content(sourceCode), -1)
			output := fmt.Sprintf(
				"%s:%d:%d:%d:%d: %s",
				path,
				c.Node.StartPoint().Row+1,
				c.Node.StartPoint().Column,
				c.Node.EndPoint().Row+1,
				c.Node.EndPoint().Column-1,
				msg,
			)

			// Should we move the print to somewhere else?
			errCount += 1
			fmt.Println(output)
		}
	}

	return errCount, nil
}
