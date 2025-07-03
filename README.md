# tint

> **T**ree-sitter powered l**int**er

- **Running**: `tint lint`
- **Config**: Sample file available in `.tint.toml.sample` (default loc is `.tint.toml`)
  For detailed instructions and examples on how to write a new rule, see [docs/writing_rules.md](docs/writing_rules.md).
- **Output format**: `filename:start-line:start-col:end-line:end-col: message`

## Configuration

The configuration file supports glob patterns for including and excluding files at both global and per-rule levels:

```toml
# Global include patterns - if specified, only files matching these patterns will be linted
include = ["**/*.go", "src/**/*.js"]

# Global exclude patterns - files matching these patterns will be excluded from linting
exclude = ["vendor/**", "node_modules/**", "**/testdata/**"]

[rules.my-rule]
language = "go"
message = "my rule message"
query = "..."
# Per-rule patterns - these are applied in addition to global patterns
include = ["**/*.go"]     # Only apply this rule to Go files
exclude = ["**/config*"]  # But not to config files

[rules.test-specific-rule]
language = "go"
message = "test-specific rule"
query = "..."
include = ["**/*_test.go"]  # Only apply to test files
```

**Pattern matching behavior:**
- **Global patterns** are applied first to determine which files to process
- **Per-rule patterns** are then applied to determine if a specific rule should run on a file
- If `include` patterns are provided (global or per-rule), only files matching at least one pattern will be processed
- If `exclude` patterns are provided (global or per-rule), files matching any pattern will be skipped
- The pattern resolution order is: Global include → Global exclude → Per-rule include → Per-rule exclude
- Patterns support `**` for recursive directory matching and `*` for single-level wildcards
- Both files and directories are checked against patterns (directories matching exclude patterns will be skipped entirely)

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
