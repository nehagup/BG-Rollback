package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"os"
)

type GitClient struct {
	gc *git.Remote
}

func (gc *GitClient) GetHead(user string, repo string) (string, error) {
	r := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/" + user + "/" + repo},
	})

	if r == nil {
		fmt.Println("Github Remote client is", r)
		return "", nil
	}

	h, err := r.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: user, // anything except an empty string
			Password: "e34543006f798c77b85d3de3bf2ea693eccf0731",
		}})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return h[0].Hash().String(), nil
}

func (gc *GitClient) GetRepoClone(user string, repo string, dir string) (*git.Repository, error) {
	return git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/" + user + "/" + repo,
		Progress: os.Stdout,
		Auth: &http.BasicAuth{
			Username: user,                                       // anything except an empty string
			Password: "e34543006f798c77b85d3de3bf2ea693eccf0731", //TODO auth
		}})
}
