package main

import (
	"BG-rollback/git"
	"BG-rollback/model"
	"BG-rollback/parser"
	"fmt"
	"log"
)

func main() {
	gc := git.New()
	p := parser.New(gc)
	
	ru := model.Repo{
		ID:   "flux-demo",
		User: "slayerjain",
		Filepath: "/clusters/demo-cluster/stack.yaml",
	}
	
	latestTag, err := p.SyncRepo(&ru)
	if err != nil {
		log.Println(err)
	}

	fmt.Print(p.CheckLatestTag(latestTag))

}