// Copyright (c) HashiCorp, Inc.

package factory

import (
	"context"
	"errors"
	"reflect"
	"terraform-provider-gitsync/internal/git/github"
	"terraform-provider-gitsync/internal/git/gitlab"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name      string
		gitURL    string
		wantHost  string
		wantOwner string
		wantRepo  string
		wantErr   error
	}{
		{
			name:      "valid without .git",
			gitURL:    "https://github.com/foo/bar",
			wantHost:  "github.com",
			wantOwner: "foo",
			wantRepo:  "bar",
			wantErr:   nil,
		},
		{
			name:      "valid https with .git",
			gitURL:    "https://github.com/foo/bar.git",
			wantHost:  "github.com",
			wantOwner: "foo",
			wantRepo:  "bar",
			wantErr:   nil,
		},
		{
			name:      "valid http scheme",
			gitURL:    "http://github.com/foo/bar.git",
			wantHost:  "github.com",
			wantOwner: "foo",
			wantRepo:  "bar",
			wantErr:   nil,
		},
		{
			name:      "get correct owner and repo with extra path",
			gitURL:    "https://gitlab.com/foo/project-1/bar.git",
			wantHost:  "gitlab.com",
			wantOwner: "foo",
			wantRepo:  "project-1/bar",
			wantErr:   nil,
		},
		{
			name:    "invalid scheme",
			gitURL:  "ssh://github.com/foo/bar.git",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "missing repo path",
			gitURL:  "https://github.com/foo",
			wantErr: ErrInvalidPath,
		},
		{
			name:    "not a valid URL",
			gitURL:  "::::",
			wantErr: ErrInvalidGitURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, owner, repo, err := parseURL(tt.gitURL)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil {
				if host != tt.wantHost {
					t.Errorf("host: got %q, want %q", host, tt.wantHost)
				}
				if owner != tt.wantOwner {
					t.Errorf("owner: got %q, want %q", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("repo: got %q, want %q", repo, tt.wantRepo)
				}
			}
		})
	}
}

func TestCreateClient(t *testing.T) {
	ctx := context.Background()

	origGitHubNewClientFunc := github.NewClientFunc
	defer func() { github.NewClientFunc = origGitHubNewClientFunc }()

	githubMockClient := &github.Client{}
	gitlabMockClient := &gitlab.Client{}

	tests := []struct {
		name           string
		url            string
		wantTypeGitHub bool
		wantTypeGitLab bool
	}{
		{
			name:           "GitHub basic client",
			url:            "https://github.com/iypetrov/terraform-provider-gitsync-e2e-test",
			wantTypeGitHub: true,
			wantTypeGitLab: false,
		},
		{
			name:           "GitLab basic client",
			url:            "https://gitlab.com/iypetrov/terraform-provider-gitsync-e2e-test",
			wantTypeGitHub: false,
			wantTypeGitLab: true,
		},
		{
			name:           "GitLab self-managed client",
			url:            "https://mycompany.com/iypetrov/terraform-provider-gitsync-e2e-test",
			wantTypeGitHub: false,
			wantTypeGitLab: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			github.NewClientFunc = func(ctx context.Context, host, owner, repo, token string) (*github.Client, error) {
				return githubMockClient, nil
			}
			gitlab.NewClientFunc = func(ctx context.Context, host, owner, repo, token string) (*gitlab.Client, error) {
				return gitlabMockClient, nil
			}

			f := NewFactory()
			client, err := f.CreateClient(ctx, tt.url, "fake-token")
			assert.NoError(t, err)

			cntTrue := 0
			for _, flag := range []bool{tt.wantTypeGitHub, tt.wantTypeGitLab} {
				if flag {
					cntTrue++
				}
			}

			exactlyOne := cntTrue == 1
			require.True(t, exactlyOne, "expected exactly one of the wantType to be true")

			if tt.wantTypeGitHub {
				assert.IsType(t, (*github.Client)(nil), client)
				assert.NotEqual(t,
					reflect.TypeOf((*gitlab.Client)(nil)),
					reflect.TypeOf(client),
					"client should not be a GitLab client",
				)
			}

			if tt.wantTypeGitLab {
				assert.IsType(t, (*gitlab.Client)(nil), client)
				assert.NotEqual(t,
					reflect.TypeOf((*github.Client)(nil)),
					reflect.TypeOf(client),
					"client should not be a GitHub client",
				)
			}
		})
	}
}
