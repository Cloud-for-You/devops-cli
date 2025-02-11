package cmd

import (
	"fmt"
	"log"
	"net/http"

	pkg "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	groupName        string
	groupDescription string
	visibility       string
)

// Create GitLab group
var CreateCmd = &cobra.Command{
	Use:                   "create",
	Short:                 "Create GitLab group",
	DisableFlagsInUseLine: true,
	Run:                   createGroup,
}

func init() {
	CreateCmd.Flags().StringVar(&groupName, "name", "", "Name of the group (required)")
	CreateCmd.Flags().StringVar(&groupDescription, "description", "", "Description of the group")
	CreateCmd.Flags().StringVar(&visibility, "visibility", "private", "Visibility of the group (private, internal, public)")

	CreateCmd.MarkFlagRequired("name")
}

func createGroup(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	result, res, err := pkg.CreateGroup(client, groupName, groupDescription, visibility)
	if err != nil {
		if res != nil && res.StatusCode == http.StatusConflict {
			fmt.Printf("Group '%s' is exists.\n", groupName)
		} else {
			fmt.Printf("Failed to create GitLab group '%s': %v\n", groupName, err)
		}
	} else {
		fmt.Printf("Group created successfully\n")
		fmt.Printf("Name: %s\n", result.Name)
		fmt.Printf("Description: %s\n", result.Description)
		fmt.Printf("Web URL: %s\n", result.WebURL)
	}
}
