package gitlab

import (
	"context"
	"terraform-provider-gitsync/internal/git"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	_ git.Client = (*Client)(nil)
)

type Client struct {
	Owner      string
	Repository string
	*gitlab.Client
}

var NewClientFunc = newClient

func newClient(ctx context.Context, owner, repo, token string) (*Client, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &Client{
		Owner:      owner,
		Repository: repo,
		Client:     client,
	}, nil
}

func (c *Client) Create(ctx context.Context, data git.ValuesYamlModel) error {
	return nil
}
