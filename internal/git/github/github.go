// Copyright (c) HashiCorp, Inc.

package github

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-gitsync/internal/git"

	"github.com/cenkalti/backoff/v5"
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

func newClient(ctx context.Context, host, owner, repo, token string) (*Client, error) {
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
		branch,
		strings.ReplaceAll(strings.ReplaceAll(path, "/", "-"), ".", "-"),
	)
}

func retryOnConflict(ctx context.Context, operation func() error) error {
	retryableOperation := func() (struct{}, error) {
		err := operation()
		if err == nil {
			return struct{}{}, nil
		}

		// Check if it's a 409 conflict error - retry these
		if ghErr, ok := err.(*github.ErrorResponse); ok && ghErr.Response.StatusCode == 409 {
			return struct{}{}, err
		}
		
		// Not a conflict error, don't retry
		return struct{}{}, backoff.Permanent(err)
	}

	_, err := backoff.Retry(ctx, retryableOperation)
	return err
}

func (c *Client) Create(ctx context.Context, data git.ValuesModel) error {
	return retryOnConflict(ctx, func() error {
		options := &github.RepositoryContentFileOptions{
			Message: github.Ptr(
				fmt.Sprintf("terraform: Create %q at branch %q", data.Path, data.Branch),
			),
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
		return err
	})
}

func (c *Client) get(ctx context.Context, path, branch string) (*github.RepositoryContent, error) {
	cnt, _, _, err := c.Repositories.GetContents(
		ctx,
		c.owner,
		c.repository,
		path,
		&github.RepositoryContentGetOptions{
			Ref: "heads/" + branch,
		},
	)
	if err != nil {
		return &github.RepositoryContent{}, err
	}

	return cnt, nil
}

func (c *Client) GetContent(ctx context.Context, path, branch string) (string, error) {
	cnt, err := c.get(ctx, path, branch)
	if err != nil {
		return "", err
	}

	if cnt == nil {
		return "", fmt.Errorf("file %q does not exist on branch %q", path, branch)
	}

	decoded, err := cnt.GetContent()
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func (c *Client) Update(ctx context.Context, data git.ValuesModel) error {
	return retryOnConflict(ctx, func() error {
		cnt, err := c.get(ctx, data.Path, data.Branch)
		if err != nil {
			return err
		}

		if cnt == nil {
			return fmt.Errorf("file %q does not exist on branch %q", data.Path, data.Branch)
		}

		sha := cnt.GetSHA()
		if sha == "" {
			return fmt.Errorf("unable to determine SHA for %q on branch %q", data.Path, data.Branch)
		}

		opts := &github.RepositoryContentFileOptions{
			Message: github.Ptr(
				fmt.Sprintf("terraform: Update %q at branch %q", data.Path, data.Branch),
			),
			Content: []byte(data.Content),
			Branch:  github.Ptr(data.Branch),
			SHA:     github.Ptr(sha),
		}

		_, _, err = c.Repositories.UpdateFile(
			ctx,
			c.owner,
			c.repository,
			data.Path,
			opts,
		)
		return err
	})
}

func (c *Client) Delete(ctx context.Context, path, branch string) error {
	return retryOnConflict(ctx, func() error {
		cnt, err := c.get(ctx, path, branch)
		if err != nil {
			return err
		}

		if cnt == nil {
			return fmt.Errorf("file %q does not exist on branch %q", path, branch)
		}

		sha := cnt.GetSHA()
		if sha == "" {
			return fmt.Errorf("unable to determine SHA for %q on branch %q", path, branch)
		}

		opts := &github.RepositoryContentFileOptions{
			Message: github.Ptr(
				fmt.Sprintf("terraform: Delete %q from branch %q", path, branch),
			),
			SHA:    github.Ptr(sha),
			Branch: github.Ptr(branch),
		}

		_, _, err = c.Repositories.DeleteFile(ctx, c.owner, c.repository, path, opts)
		return err
	})
}

func (c *Client) Owner() string {
	return c.owner
}

func (c *Client) Repository() string {
	return c.repository
}
