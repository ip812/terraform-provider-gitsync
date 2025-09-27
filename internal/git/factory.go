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
	Create(ctx context.Context, data ValuesYamlModel) error
	Owner() string
	Repository() string
}
