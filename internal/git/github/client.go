package github

import (
	"context"
	"terraform-provider-gitsync/internal/git"

	"github.com/google/go-github/v75/github"
	gg "github.com/google/go-github/v75/github"
	"golang.org/x/oauth2"
)

var (
	_ git.Client = (*Client)(nil)
)

type Client struct {
	Owner      string
	Repository string
	*gg.Client
}

func New(ctx context.Context, owner, repo, token string) (*Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gg.NewClient(tc)
	return &Client{
		Owner:      owner,
		Repository: repo,
		Client:     client,
	}, nil
}

func (c *Client) Create(ctx context.Context, data git.ValuesYamlModel) error {
	options := &github.RepositoryContentFileOptions{
		Message: github.Ptr("Update values.yaml from Terraform"),
		Content: []byte(data.Content),
		Branch:  github.Ptr(data.Branch),
	}

	_, _, err := c.Repositories.CreateFile(
		ctx,
		c.Owner,
		c.Repository,
		data.Path,
		options,
	)
	if err != nil {
		return err
	}

	return nil
}
