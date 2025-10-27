# Test Fixtures

This directory contains test fixture files used by the validators package tests.

## Directory Structure

```
fixtures/
└── formats/
    ├── valid_*.yaml      - Valid YAML test files
    ├── valid_*.json      - Valid JSON test files
    ├── invalid_*.yaml    - Invalid YAML test files
    ├── invalid_*.json    - Invalid JSON test files
    └── empty.txt         - Empty file for testing empty content validation
```

## File Naming Convention

- `valid_*` - Files that should pass validation
- `invalid_*` - Files that should fail validation
- Extensions should match the format being tested (`.yaml`, `.json`, etc.)

## Usage

These fixtures are loaded by the `loadFixture()` helper function in `validators_test.go`:

```go
content := loadFixture(t, "valid_simple.yaml")
err := ValidateYAML(content)
```

## Adding New Fixtures

1. Create a new file in `fixtures/formats/`
2. Follow the naming convention (`valid_*` or `invalid_*`)
3. Add a corresponding test case in `validators_test.go`
4. Update this README if adding a new category of fixtures

## Fixture Categories

### YAML Fixtures

**Valid:**
- `valid_simple.yaml` - Simple key-value pairs
- `valid_nested.yaml` - Nested structures
- `valid_arrays.yaml` - YAML with arrays
- `valid_special_chars.yaml` - Special characters in values
- `valid_plain_string.yaml` - Plain string (valid YAML)
- `valid_number.yaml` - Number (valid YAML)

**Invalid:**
- `invalid_bad_indentation.yaml` - Incorrect indentation
- `invalid_unclosed_bracket.yaml` - Unclosed array bracket
- `invalid_mixed_tabs.yaml` - Mixed tabs and spaces

### JSON Fixtures

**Valid:**
- `valid_simple.json` - Simple JSON object
- `valid_nested.json` - Nested objects and arrays
- `valid_array.json` - JSON array
- `valid_minimal_object.json` - Empty object `{}`
- `valid_minimal_array.json` - Empty array `[]`
- `valid_null.json` - JSON with null values
- `valid_special_chars.json` - Special characters
- `valid_escaped_chars.json` - Escaped characters

**Invalid:**
- `invalid_missing_brace.json` - Missing closing brace
- `invalid_trailing_comma.json` - Trailing comma
- `invalid_single_quotes.json` - Single quotes (not valid JSON)
- `invalid_unquoted_keys.json` - Unquoted keys
- `invalid_comments.json` - Comments (not allowed in JSON)
- `invalid_plain_string.json` - Plain string (not valid JSON)

### Common Fixtures

- `empty.txt` - Empty file for testing empty content validation

