package sync

import (
	"os"

	"github.com/go-git/go-billy/v6/memfs"
	"github.com/go-git/go-git/v6"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/storage/memory"
)

type Client struct {
	Repo *git.Repository
}

func NewClient(repository, branch, token *string) (*Client, error) {
	fs := memfs.New()
	storer := memory.NewStorage()

	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: *repository,
		Auth: &githttp.BasicAuth{
			Username: "Terraform Syncronizer",
			Password: *token,
		},
		Depth:        1,
		Tags:         git.NoTags,
		SingleBranch: true,
		Progress:     os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		Repo: repo,
	}, nil
}
