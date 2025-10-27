// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFixture(t *testing.T, filename string) string {
	t.Helper()
	path := filepath.Join("..", "..", "fixtures", "formats", filename)
	content, err := os.ReadFile(path)
	require.NoError(t, err, "failed to load fixture file: %s", filename)
	return string(content)
}

func TestValidateYAML(t *testing.T) {
	tests := []struct {
		name        string
		fixture     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid simple yaml",
			fixture:     "valid_simple.yaml",
			expectError: false,
		},
		{
			name:        "valid nested yaml",
			fixture:     "valid_nested.yaml",
			expectError: false,
		},
		{
			name:        "valid yaml with arrays",
			fixture:     "valid_arrays.yaml",
			expectError: false,
		},
		{
			name:        "empty content",
			fixture:     "empty.txt",
			expectError: true,
			errorMsg:    "content cannot be empty",
		},
		{
			name:        "invalid yaml - bad indentation",
			fixture:     "invalid_bad_indentation.yaml",
			expectError: true,
			errorMsg:    "failed to parse YAML",
		},
		{
			name:        "invalid yaml - unclosed bracket",
			fixture:     "invalid_unclosed_bracket.yaml",
			expectError: true,
			errorMsg:    "failed to parse YAML",
		},
		{
			name:        "invalid yaml - mixed tabs and spaces",
			fixture:     "invalid_mixed_tabs.yaml",
			expectError: true,
			errorMsg:    "failed to parse YAML",
		},
		{
			name:        "plain string (still valid yaml)",
			fixture:     "valid_plain_string.yaml",
			expectError: false,
		},
		{
			name:        "number (valid yaml)",
			fixture:     "valid_number.yaml",
			expectError: false,
		},
		{
			name:        "yaml with special characters",
			fixture:     "valid_special_chars.yaml",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := loadFixture(t, tt.fixture)
			err := ValidateYAML(content)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				
				validationErr, ok := err.(*ValidationError)
				require.True(t, ok, "error should be a ValidationError")
				assert.Equal(t, "yaml", validationErr.Type)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateJSON(t *testing.T) {
	tests := []struct {
		name        string
		fixture     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid simple json",
			fixture:     "valid_simple.json",
			expectError: false,
		},
		{
			name:        "valid nested json",
			fixture:     "valid_nested.json",
			expectError: false,
		},
		{
			name:        "valid json array",
			fixture:     "valid_array.json",
			expectError: false,
		},
		{
			name:        "valid minimal json object",
			fixture:     "valid_minimal_object.json",
			expectError: false,
		},
		{
			name:        "valid minimal json array",
			fixture:     "valid_minimal_array.json",
			expectError: false,
		},
		{
			name:        "valid json with null",
			fixture:     "valid_null.json",
			expectError: false,
		},
		{
			name:        "empty content",
			fixture:     "empty.txt",
			expectError: true,
			errorMsg:    "content cannot be empty",
		},
		{
			name:        "invalid json - missing closing brace",
			fixture:     "invalid_missing_brace.json",
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name:        "invalid json - trailing comma",
			fixture:     "invalid_trailing_comma.json",
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name:        "invalid json - single quotes",
			fixture:     "invalid_single_quotes.json",
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name:        "invalid json - unquoted keys",
			fixture:     "invalid_unquoted_keys.json",
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name:        "invalid json - comments not allowed",
			fixture:     "invalid_comments.json",
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name:        "plain string (not valid json)",
			fixture:     "invalid_plain_string.json",
			expectError: true,
			errorMsg:    "failed to parse JSON",
		},
		{
			name:        "valid json with special characters",
			fixture:     "valid_special_chars.json",
			expectError: false,
		},
		{
			name:        "valid json with escaped characters",
			fixture:     "valid_escaped_chars.json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := loadFixture(t, tt.fixture)
			err := ValidateJSON(content)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				
				validationErr, ok := err.(*ValidationError)
				require.True(t, ok, "error should be a ValidationError")
				assert.Equal(t, "json", validationErr.Type)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "error with underlying error",
			err: &ValidationError{
				Type:    "json",
				Message: "parsing failed",
				Err:     assert.AnError,
			},
			expected: "invalid json: parsing failed: assert.AnError general error for testing",
		},
		{
			name: "error without underlying error",
			err: &ValidationError{
				Type:    "yaml",
				Message: "content cannot be empty",
			},
			expected: "invalid yaml: content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

