# devops-cli

## Create Project
```shell
GITLAB_URL="https://gitlab.pmb.cz"
GITLAB_TOKEN="pat-fUScVcKtosoqQ1KrnkTQ"
GITLAB_NAMESPACE="480"
MAINTAINER_GROUP="GIT_CLI_Maintainer_G"
DEVELOPER_GROUP="GIT_CLI_Developer_G"
PROJECT_NAME="cli"

$ go run main.go gitlab-ce project create --gitlabUrl=${GITLAB_URL} --gitlabToken=${GITLAB_TOKEN} --namespace=${GITLAB_NAMESPACE} --maintainerGroup=${MAINTAINER_GROUP} --developerGroup=${DEVELOPER_GROUP} --name=${PROJECT_NAME}
```