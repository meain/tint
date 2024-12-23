# tint

> **T**ree-sitter powered l**int**er

- **Running**: `tint lint`
- **Config**: Sample file available in `.tint.lint.sample` (default loc is `.tint.lint`)
- **Output format**: `filename:start-line:start-col:end-line:end-col: message`

### Installation

```
go install github.com/meain/tint@latest
```

### Example output:

```
config.go:72:4:72:9: do not use "" to check for empty string for 'config'
lint.go:83:7:83:7: do not use trailing comma for args
main.go:122:16:122:16: do not use trailing comma for args
main.go:131:39:131:39: do not use trailing comma for args
```

### Help

```
Usage: tint <command> [flags]

Flags:
  -h, --help             Show context-sensitive help.
      --config=STRING    Path to config file

Commands:
  lint [<files> ...] [flags]
    Lint files or folders

  validate-config [flags]
    Validate config file

Run "tint <command> --help" for more information on a command.
```
