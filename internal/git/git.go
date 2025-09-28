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
	GetContent(ctx context.Context, path, branch string) (string, error)
	Update(ctx context.Context, data ValuesYamlModel) error
	Delete(ctx context.Context, path, branch string) error
	Owner() string
	Repository() string
}
