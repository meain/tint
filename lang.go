package main

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type Language string

const (
	LanguageGo Language = "go"
)

type LangInfo struct {
	Names      []string
	Extensions []string
	TSLang     *sitter.Language
}

var LanguageMap = map[Language]LangInfo{
	LanguageGo: {Names: []string{"go", "golang"}, Extensions: []string{"go"}, TSLang: golang.GetLanguage()},
}
