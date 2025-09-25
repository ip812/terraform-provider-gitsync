package factory

import (
	"errors"
	"testing"
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
