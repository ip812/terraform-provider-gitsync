package factory

import (
	"errors"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name      string
		gitURL    string
		wantOwner string
		wantRepo  string
		wantErr   error
	}{
		{
			name:      "valid https with .git",
			gitURL:    "https://github.com/iypetrov/terraform-provider-gitsync-e2e-test.git",
			wantOwner: "iypetrov",
			wantRepo:  "terraform-provider-gitsync-e2e-test",
			wantErr:   nil,
		},
		{
			name:      "valid https without .git",
			gitURL:    "https://github.com/golang/go",
			wantOwner: "golang",
			wantRepo:  "go",
			wantErr:   nil,
		},
		{
			name:      "valid http scheme",
			gitURL:    "http://github.com/foo/bar.git",
			wantOwner: "foo",
			wantRepo:  "bar",
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
			owner, repo, err := parseURL(tt.gitURL)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got err %v, want %v", err, tt.wantErr)
			}
			if err == nil {
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
