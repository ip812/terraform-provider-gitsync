package github

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-gitsync/internal/git"

	"github.com/google/go-github/v75/github"
	"golang.org/x/oauth2"
)

var (
	_ git.Client = (*Client)(nil)
)

type Client struct {
	owner      string
	repository string
	*github.Client
}

var NewClientFunc = newClient

func newClient(ctx context.Context, owner, repo, token string) (*Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &Client{
		owner:      owner,
		repository: repo,
		Client:     github.NewClient(tc),
	}, nil
}

func (c *Client) GetID(branch, path string) string {
	return fmt.Sprintf(
		"github-%s-%s-%s-%s",
		c.owner,
		c.repository,
		strings.ReplaceAll(branch, "/", "-"),
		strings.ReplaceAll(path, "/", "-"),
	)
}

func (c *Client) Create(ctx context.Context, data git.ValuesYamlModel) error {
	options := &github.RepositoryContentFileOptions{
		Message: github.Ptr("Update values.yaml from Terraform"),
		Content: []byte(data.Content),
		Branch:  github.Ptr(data.Branch),
	}

	_, _, err := c.Repositories.CreateFile(
		ctx,
		c.owner,
		c.repository,
		data.Path,
		options,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetContent(ctx context.Context, id, path string) (string, error) {
	cnt, _, _, err := c.Repositories.GetContents(
		ctx,
		c.owner,
		c.repository,
		path,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil {
		return "", err
	}

	if cnt == nil {
		return "", nil
	}

	decoded, err := cnt.GetContent()
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func (c *Client) Owner() string {
	return c.owner
}

func (c *Client) Repository() string {
	return c.repository
}
