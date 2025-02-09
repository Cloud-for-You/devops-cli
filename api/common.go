package api

import "encoding/json"

const Domain = "cfy.cz"
const GitLabAPIVersion = "gitlab-ce.devops-cli.cfy.cz/v1alpha1"

type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

func Validate(data []byte, v any) bool {
	err := json.Unmarshal(data, v)
	return err == nil
}
