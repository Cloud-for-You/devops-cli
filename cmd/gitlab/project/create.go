package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	v1alpha1 "github.com/Cloud-for-You/devops-cli/api/v1alpha1/gitlab-ce"
	common "github.com/Cloud-for-You/devops-cli/pkg"
	pkg "github.com/Cloud-for-You/devops-cli/pkg/gitlab"
	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var (
	projectName         string
	projectDescription  string
	namespaceID         int
	visibility          string
	maintainerGroupName string
  developerGroupName  string	
	importURL					  string
	configFile					string
)

// Create GitLab repository
var CreateCmd = &cobra.Command{
	Use:                   "create",
	Short:                 "Create GitLab repository",
	DisableFlagsInUseLine: true,
	Run:                   createProject,
}

func init() {
	CreateCmd.Flags().StringVar(&projectName, "name", "", "Name of the repository (required)")
	CreateCmd.Flags().StringVar(&projectDescription, "description", "", "Description of the repository")
	CreateCmd.Flags().IntVar(&namespaceID, "namespaceID", 0, "Namespace ID under which the repository will be created")
	CreateCmd.Flags().StringVar(&visibility, "visibility", "private", "Visibility of the repository (private, internal, public)")
	CreateCmd.Flags().StringVar(&maintainerGroupName, "maintainerGroup", "", "Group containing maintainers")
	CreateCmd.Flags().StringVar(&developerGroupName, "developerGroup", "", "Group containing developers")
	CreateCmd.Flags().StringVar(&importURL, "importURL", "", "URL of the repository to import")
	CreateCmd.Flags().StringVarP(&configFile, "configFile", "f", "", "GitLab Project config file from options: json, yaml")

	//CreateCmd.MarkFlagRequired("name")
}

func createProject(cmd *cobra.Command, args []string) {
	gitlabToken, _ := cmd.Flags().GetString("gitlabToken")
	gitlabUrl, _ := cmd.Flags().GetString("gitlabUrl")

	// Nadefinujeme options pro vytvoreni gitlab projektu podle specifik.
	// Nejdrive nacteme konfiguraci ze souboru pokud byl soubor predan
	// a nasledne pretizime hodnoty z dodanych flags
	projectSpec := pkg.ProjectSpec{}

	// Nacteme konfiguraci ze souboru nebo URL pokud jsme jej definovali v commandLine
	if configFile != "" {
		data, err := common.LoadConfigFile(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		// Validujeme konfiguraci v souboru na strukturu GitLabProject
  	dataBytes, err := json.Marshal(data)
  	if err != nil {
   		log.Fatalf("Error converting config data to JSON: %v", err)
  	}
	
	  var gitLabProject v1alpha1.GitLabProject
			err = json.Unmarshal(dataBytes, &gitLabProject)
			if err != nil {
				log.Fatalf("Error unmarshalling config data: %v", err)
			}
			projectSpec = gitLabProject.Spec
			projectSpec.Name = &gitLabProject.Metadata.Name
	}

	// Pretizime konfiguraci z commandLine parametru
	if cmd.Flags().Changed("name") {
		projectSpec.Name = &projectName
	}	

	if cmd.Flags().Changed("namespaceID") {
		projectSpec.NamespaceID = &namespaceID
	}

	if cmd.Flags().Changed("description") {
		projectSpec.Description = &projectDescription
	}

	if cmd.Flags().Changed("importURL") {
		projectSpec.ImportURL = &importURL
	}
	
	if cmd.Flags().Changed("visibility") {
		projectSpec.Visibility = (*gitlab.VisibilityValue)(&visibility)
	}

	// Provedeme pripojeni ke GitLabu a zalozeni projektu
	if gitlabToken == "" || gitlabUrl == "" {
		log.Fatalf("Gitlab token and URL must be provided using the persistent flags --gitlabToken and --gitlabUrl")
	}

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabUrl))
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	result, res, err := pkg.CreateProject(
		client,
		&projectSpec,
		&maintainerGroupName,
		&developerGroupName,
	) 
	if err != nil {
		if res != nil && res.StatusCode == http.StatusConflict {
			fmt.Printf("Project '%s' is exists.\n", projectName)
		} else {
			fmt.Printf("Failed to create GitLab project '%s': %v\n", projectName, err)
		}
	} else {
		fmt.Printf("Project created successfully\n")
		fmt.Printf("Name: %s\n", result.Name)
		fmt.Printf("Description: %s\n", result.Description)
		fmt.Printf("Web URL: %s\n", result.WebURL)
	}
}
