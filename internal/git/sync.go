package git

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v75/github"
	"golang.org/x/oauth2"
)

type Client struct {
	GitHubOwner      string
	GitHubRepository string
	GitHubClient     *github.Client
}

func parseGitHubURL(repoURL string) (owner string, repo string, err error) {
	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected repo URL format")
	}

	owner = parts[0]
	repo = strings.TrimSuffix(parts[1], ".git")
	return owner, repo, nil
}

func NewClient(ctx context.Context, url, token string) (*Client, error) {
	owner, repo, err := parseGitHubURL(url)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &Client{
		GitHubOwner:      owner,
		GitHubRepository: repo,
		GitHubClient:     client,
	}, nil
}
