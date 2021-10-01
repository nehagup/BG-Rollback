package model

import (
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta1"
)
type Repo struct {
	ID       string `json:"repo"`
	User     string `json:"user"`
	Filepath string `json:"filepath"`
}

type SourceCRDs struct {
	GitRepository  map[string]sourcev1.GitRepository
	Bucket         map[string]sourcev1.Bucket
	HelmChart      map[string]sourcev1.HelmChart
	HelmRepository map[string]sourcev1.HelmRepository
}

