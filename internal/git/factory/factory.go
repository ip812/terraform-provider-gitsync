// Copyright (c) HashiCorp, Inc.

package factory

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"terraform-provider-gitsync/internal/git"
	"terraform-provider-gitsync/internal/git/github"
	"terraform-provider-gitsync/internal/git/gitlab"
)

var (
	ErrInvalidGitURL     = fmt.Errorf("invalid git URL")
	ErrUnsupportedScheme = fmt.Errorf("unsupported URL scheme")
	ErrInvalidPath       = fmt.Errorf("invalid git URL path, expected format: <host>/<owner>/<repo>")
)

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

// Use GitHub client only for github.com; all other hosts fall back to GitLab client.
// This covers both gitlab.com and self-hosted GitLab instances with custom domains.
func (f *Factory) CreateClient(ctx context.Context, url, token string) (git.Client, error) {
	host, owner, repo, err := parseURL(url)
	if err != nil {
		return nil, err
	}

	if host == "github.com" {
		client, err := github.NewClientFunc(ctx, host, owner, repo, token)
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	client, err := gitlab.NewClientFunc(ctx, host, owner, repo, token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func parseURL(gitURL string) (host, owner, repo string, err error) {
	u, err := url.ParseRequestURI(gitURL)
	if err != nil {
		return "", "", "", ErrInvalidGitURL
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return "", "", "", ErrUnsupportedScheme
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", "", ErrInvalidPath
	}

	host = u.Host
	owner = parts[0]
	repo = path.Join(parts[1:]...)
	repo = strings.TrimSuffix(repo, ".git")

	return host, owner, repo, nil
}
