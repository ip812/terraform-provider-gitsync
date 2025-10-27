// Copyright (c) HashiCorp, Inc.

package gitlab

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-gitsync/internal/git"

	"github.com/cenkalti/backoff/v5"
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

func newClient(ctx context.Context, host, owner, repo, token string) (*Client, error) {
	client, err := gitlab.NewClient(
		token,
		gitlab.WithBaseURL(fmt.Sprintf("https://%s/api/v4", host)),
	)
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

func (c *Client) projectPath() string {
	return fmt.Sprintf("%s/%s", c.owner, c.repository)
}

func retryOnConflict(ctx context.Context, operation func() error) error {
	retryableOperation := func() (struct{}, error) {
		err := operation()
		if err == nil {
			return struct{}{}, nil
		}

		if glErr, ok := err.(*gitlab.ErrorResponse); ok {
			statusCode := glErr.Response.StatusCode
			if statusCode == http.StatusConflict || statusCode == http.StatusBadRequest {
				return struct{}{}, err
			}
		}

		return struct{}{}, backoff.Permanent(err)
	}

	_, err := backoff.Retry(ctx, retryableOperation)
	return err
}

func (c *Client) Create(ctx context.Context, data git.ValuesModel) error {
	return retryOnConflict(ctx, func() error {
		msg := fmt.Sprintf("terraform: Create %q at branch %q", data.Path, data.Branch)
		opts := &gitlab.CreateFileOptions{
			Branch:        gitlab.Ptr(data.Branch),
			Content:       gitlab.Ptr(data.Content),
			CommitMessage: gitlab.Ptr(msg),
		}

		_, _, err := c.RepositoryFiles.CreateFile(c.projectPath(), data.Path, opts, gitlab.WithContext(ctx))
		return err
	})
}

func (c *Client) get(ctx context.Context, path, branch string) (*gitlab.File, error) {
	file, _, err := c.RepositoryFiles.GetFile(
		c.projectPath(),
		path,
		&gitlab.GetFileOptions{Ref: gitlab.Ptr(branch)},
		gitlab.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *Client) GetContent(ctx context.Context, path, branch string) (string, error) {
	file, err := c.get(ctx, path, branch)
	if err != nil {
		return "", err
	}
	if file == nil {
		return "", fmt.Errorf("file %q does not exist on branch %q", path, branch)
	}

	decoded, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func (c *Client) Update(ctx context.Context, data git.ValuesModel) error {
	return retryOnConflict(ctx, func() error {
		file, err := c.get(ctx, data.Path, data.Branch)
		if err != nil {
			return err
		}
		if file == nil {
			return fmt.Errorf("file %q does not exist on branch %q", data.Path, data.Branch)
		}

		msg := fmt.Sprintf("terraform: Update %q at branch %q", data.Path, data.Branch)
		opts := &gitlab.UpdateFileOptions{
			Branch:        gitlab.Ptr(data.Branch),
			Content:       gitlab.Ptr(data.Content),
			CommitMessage: gitlab.Ptr(msg),
			LastCommitID:  gitlab.Ptr(file.LastCommitID),
		}

		_, _, err = c.RepositoryFiles.UpdateFile(
			c.projectPath(),
			data.Path,
			opts,
			gitlab.WithContext(ctx),
		)
		return err
	})
}

func (c *Client) Delete(ctx context.Context, path, branch string) error {
	return retryOnConflict(ctx, func() error {
		file, err := c.get(ctx, path, branch)
		if err != nil {
			return err
		}
		if file == nil {
			return fmt.Errorf("file %q does not exist on branch %q", path, branch)
		}

		msg := fmt.Sprintf("terraform: Delete %q from branch %q", path, branch)
		opts := &gitlab.DeleteFileOptions{
			Branch:        gitlab.Ptr(branch),
			CommitMessage: gitlab.Ptr(msg),
			LastCommitID:  gitlab.Ptr(file.LastCommitID),
		}

		_, err = c.RepositoryFiles.DeleteFile(
			c.projectPath(),
			path,
			opts,
			gitlab.WithContext(ctx),
		)
		return err
	})
}

func (c *Client) Owner() string {
	return c.owner
}

func (c *Client) Repository() string {
	return c.repository
}
