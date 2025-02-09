package api

import (
	api "github.com/Cloud-for-You/devops-cli/api"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Project struct {
	APIVersion string       `json:"apiVersion" yaml:"apiVersion"`
	Kind       string       `json:"kind" yaml:"kind"`
	Metadata   api.Metadata `json:"metadata" yaml:"metadata"`
	Spec       ProjectSpec  `json:"spec" yaml:"spec"`
}

type ProjectSpec gitlab.CreateProjectOptions
