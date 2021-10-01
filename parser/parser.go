package parser

import (
	"BG-rollback/model"
	"bytes"
	"fmt"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	yaml2 "gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
)

func New( gc GitClient) *Parser {
	return &Parser{

		gc:  gc,
	}
}

func (p *Parser) SyncRepo(r *model.Repo) (latesttag string, err error) {
	// Tempdir to clone the repository
	dir, err := ioutil.TempDir("/tmp/", "flux-")
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println("Removing temp dir failed")
		}
	}(dir)

	_, err = p.gc.GetRepoClone(r.User, r.ID, dir)
	if err != nil {
		return "",err
	}

	// Parse apps from the locally cloned repo where key is app name
	data, err := ioutil.ReadFile(dir+r.Filepath)
	if err != nil {
		fmt.Println("File reading error", err)
		log.Println("Reading filepath failed", r.Filepath)
	}

	temp, err := p.splitYAML(data)
	for _, v := range temp {
		 gr := sourcev1.GitRepository{}
		//var d map[string]*model.Destination

		err = yaml.Unmarshal(v, &gr)
		if err != nil {
			log.Println("error: %v", err)
			continue
		}

		latestTag := gr.Spec.Reference.Tag
		return latestTag, nil

		break
	}

	return "",nil
}

func (p *Parser) splitYAML(resources []byte) ([][]byte, error) {
	dec := yaml2.NewDecoder(bytes.NewReader(resources))
	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml2.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}

func (p *Parser) CheckLatestTag(tag string) bool {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/slayerjain/demo-stack",
	})
	if err != nil {
		panic(err)
	}

	t, err := GetLatestTagFromRepository(r)
	if err !=nil {
		panic(err)
	}
	return strings.Contains(t, tag)
}

func GetLatestTagFromRepository(repository *git.Repository) (string, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return "", err
	}

	var latestTagCommit *object.Commit
	var latestTagName string
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		revision := plumbing.Revision(tagRef.Name().String())
		tagCommitHash, err := repository.ResolveRevision(revision)
		if err != nil {
			return err
		}

		commit, err := repository.CommitObject(*tagCommitHash)
		if err != nil {
			return err
		}

		if latestTagCommit == nil {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		if commit.Committer.When.After(latestTagCommit.Committer.When) {
			latestTagCommit = commit
			latestTagName = tagRef.Name().String()
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return latestTagName, nil
}


