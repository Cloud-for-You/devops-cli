package v1alpha1

import (
	api "github.com/Cloud-for-You/devops-cli/apis/v1alpha1"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitLabProject struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Metadata   api.Metadata `json:"metadata" yaml:"metadata"`
	Spec       gitlab.CreateProjectOptions `json:"spec" yaml:"spec"`
}
