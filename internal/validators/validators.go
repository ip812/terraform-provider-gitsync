// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validators

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type ValidationError struct {
	Type    string
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("invalid %s: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("invalid %s: %s", e.Type, e.Message)
}

func ValidateYAML(content string) error {
	if content == "" {
		return &ValidationError{
			Type:    "yaml",
			Message: "content cannot be empty",
		}
	}

	var data interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return &ValidationError{
			Type:    "yaml",
			Message: "failed to parse YAML content",
			Err:     err,
		}
	}

	return nil
}

func ValidateJSON(content string) error {
	if content == "" {
		return &ValidationError{
			Type:    "json",
			Message: "content cannot be empty",
		}
	}

	var data interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return &ValidationError{
			Type:    "json",
			Message: "failed to parse JSON content",
			Err:     err,
		}
	}

	return nil
}
