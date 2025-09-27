package git

import (
	"context"
)

type ValuesYamlModel struct {
	Path    string
	Branch  string
	Content string
}

type Client interface {
	GetID(branch, path string) string
	Create(ctx context.Context, data ValuesYamlModel) error
	GetContent(ctx context.Context, id, path string) (string, error)
	Owner() string
	Repository() string
}
