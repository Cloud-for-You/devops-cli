package v1alpha1

import (
	api "github.com/Cloud-for-You/devops-cli/api"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitLabGroup struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   api.Metadata `json:"metadata" yaml:"metadata"`
	Spec       gitlab.CreateGroupOptions `json:"spec" yaml:"spec"`
}
