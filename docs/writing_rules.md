# Writing New Rules

This document provides detailed instructions on how to write new rules for the tint linter.

## Rule Structure

Rules are defined in the configuration file (e.g., `.tint.toml`) under the `[rules]` section. Each rule is defined in its own subsection and must have the following fields:

- `language`: Specifies the language the rule applies to. Currently, only "go" is supported.
- `message`: The message that will be displayed when a violation of the rule is found. You can use `{}` as a placeholder in the message, which will be replaced with the text captured by the `@region` capture in the query.
- `query`: A tree-sitter query that identifies the code patterns that violate the rule.


### Example 1: Disallowing Trailing Commas in Argument Lists

```toml
[rules.no-trailing-comma]
language = "go"
message = "Do not use trailing commas in argument lists"
query = '''(argument_list "," @region)'''
```

This rule targets Go code and flags any `argument_list` that ends with a comma.

- `argument_list`: Matches an argument list node.
- `","`: Matches a comma token.
- `@region`: This capture name marks the specific comma that violates the rule. This is used to determine the location of the error and to fill the `{}` placeholder in the message (though this message doesn't use it).

### Example 2: Disallowing String Comparisons for Empty Strings

```toml
[rules.no-string-comparison]
language = "go"
message = "Do not use empty string literals (`""`) for comparison; use `len(s) == 0` instead"
query = '''
((binary_expression
  (identifier) @region ["==" "!="] (string_literal (string_content)))
 (#match? @region "^\"\"$"))
'''
```

This rule flags comparisons with empty string literals.

- `binary_expression`: Matches a binary expression (e.g., `a == ""`).
- `(identifier) @region`: Matches an identifier on either side of the operator.
- `["==" "!="]`: Matches either the `==` or `!=` operator.
- `(string_literal (string_content))`: Matches a string literal.
- `(#match? @region "^\"\"$")`: A predicate that ensures the matched string literal is an empty string.

### Example 3: Enforcing Standard Receiver Names for Test Suites

```toml
[rules.standard-name-for-test-suite]
language = "go"
message = "Use 'suite' as the receiver name for test suites"
query = '''
((method_declaration
  (parameter_list
   (parameter_declaration (identifier) @name
                          (pointer_type) @type))
  (field_identifier) @region)
 (#not-eq? @name "suite")
 (#match? @type "Suite$")
 (#match? @region "^Test"))
'''
```

This rule enforces the convention of using `suite` as the receiver name for methods belonging to types with names ending in "Suite" that are also test functions (names starting with "Test").

- `method_declaration`: Matches a method declaration.
- `parameter_list`: Matches the parameter list of the method.
- `parameter_declaration`: Matches a parameter declaration.
- `(identifier) @name`: Captures the identifier (receiver name).
- `(pointer_type) @type`: Captures the receiver type.
- `field_identifier`: Matches the method name.
- `@region`: Captures the entire method name.
- `(#not-eq? @name "suite")`: Ensures the receiver name is not "suite".
- `(#match? @type "Suite$")`: Ensures the receiver type ends with "Suite".
- `(#match? @region "^Test")`: Ensures the method name starts with "Test".

## Testing Your Rules

After adding a new rule to the configuration file, you can test it by running the `tint lint` command on a file that you expect to violate the rule.

## Learning More About Tree-sitter Queries

- [Tree-sitter Query Documentation](https://tree-sitter.github.io/tree-sitter/syntax-highlighting#queries)
- [Tree-sitter Playground](https://tree-sitter.github.io/tree-sitter/playground) - A useful tool for experimenting with queries on different code snippets.

Remember to define clear messages and precise queries to make your rules effective and helpful.
