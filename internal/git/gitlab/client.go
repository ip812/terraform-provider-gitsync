// Copyright (c) HashiCorp, Inc.

package gitlab

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-gitsync/internal/git"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	_ git.Client = (*Client)(nil)
)

type Client struct {
	owner      string
	repository string
	*gitlab.Client
}

var NewClientFunc = newClient

func newClient(ctx context.Context, owner, repo, token string) (*Client, error) {
	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return &Client{
		owner:      owner,
		repository: repo,
		Client:     client,
	}, nil
}

func (c *Client) GetID(branch, path string) string {
	return fmt.Sprintf(
		"gitlab-%s-%s-%s-%s",
		c.owner,
		c.repository,
		branch,
		strings.ReplaceAll(strings.ReplaceAll(path, "/", "-"), ".", "-"),
	)
}

func (c *Client) Create(ctx context.Context, data git.ValuesYamlModel) error {
	return nil
}

func (c *Client) GetContent(ctx context.Context, path, branch string) (string, error) {
	return "", nil
}

func (c *Client) Update(ctx context.Context, data git.ValuesYamlModel) error {
	return nil
}

func (c *Client) Delete(ctx context.Context, path, branch string) error {
	return nil
}

func (c *Client) Owner() string {
	return c.owner
}

func (c *Client) Repository() string {
	return c.repository
}
