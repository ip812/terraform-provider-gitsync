package factory

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"terraform-provider-gitsync/internal/git"
	"terraform-provider-gitsync/internal/git/github"
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

func (f *Factory) CreateClient(ctx context.Context, url, token string) (git.Client, error) {
	owner, repo, err := parseURL(url)
	if err != nil {
		return nil, err
	}

	client, err := github.New(ctx, owner, repo, token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func parseURL(gitURL string) (owner, repo string, err error) {
	u, err := url.ParseRequestURI(gitURL)
	if err != nil {
		return "", "", ErrInvalidGitURL
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return "", "", ErrUnsupportedScheme
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", ErrInvalidPath
	}

	owner = parts[0]
	repo = parts[1]
	repo = strings.TrimSuffix(repo, ".git")

	return owner, repo, nil
}
