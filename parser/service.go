package parser

import "github.com/go-git/go-git/v5"

type Parser struct {
	gc  GitClient
}

type GitClient interface {
	GetHead(user string, repo string) (string, error)
	GetRepoClone(user string, repo string, dir string) (*git.Repository, error)
}