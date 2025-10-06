// Copyright (c) HashiCorp, Inc.

package git

import (
	"context"
)

type ValuesModel struct {
	Path    string
	Branch  string
	Content string
}

type Client interface {
	GetID(branch, path string) string
	Create(ctx context.Context, data ValuesModel) error
	GetContent(ctx context.Context, path, branch string) (string, error)
	Update(ctx context.Context, data ValuesModel) error
	Delete(ctx context.Context, path, branch string) error
	Owner() string
	Repository() string
}
