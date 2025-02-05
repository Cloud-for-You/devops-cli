package gitlab

import (
	"fmt"
	"log"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func ListProjects(client *gitlab.Client) ([]*gitlab.Project, error) {
	var allProjects []*gitlab.Project
	page := 1
	perPage := 20

	for {
		options := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: perPage,
			},
		}

		projects, res, err := client.Projects.ListProjects(options)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve projects: %v", err)
		}

		allProjects = append(allProjects, projects...)

		if res.CurrentPage >= res.TotalPages {
			break
		}

		page++
	}

	return allProjects, nil
}

func CreateProject(client *gitlab.Client, projectName string, namespaceID int, projectDescription string, visibility string, maintainerGroupName *string, developerGroupName *string) (*gitlab.Project, *gitlab.Response, error) {

  projectOptions := &gitlab.CreateProjectOptions{
	  Name:        gitlab.Ptr(projectName),
		Path:        gitlab.Ptr(projectName),
		Description: gitlab.Ptr(projectDescription),
		NamespaceID: gitlab.Ptr(namespaceID),
		Visibility:  gitlab.Ptr(gitlab.VisibilityValue(visibility)),
	}

	/*
	maintainerGroup, res, err := CreateGroup(client, *maintainerGroupName, "", "private")
	if err != nil {
		if res != nil && res.StatusCode == http.StatusConflict {
			fmt.Printf("Group '%s' is exists.\n", *maintainerGroupName)
		} else {
			fmt.Printf("Failed to create GitLab group '%s': %v\n", *maintainerGroupName, err)
		}
	}

	developerGroup, res, err := CreateGroup(client, *developerGroupName, "", "private")
	if err != nil {
		if res != nil && res.StatusCode == http.StatusConflict {
			fmt.Printf("Group '%s' is exists.\n", *developerGroupName)
		} else {
			fmt.Printf("Failed to create GitLab group '%s': %v\n", *developerGroupName, err)
		}
	}
	*/

	project, res, err := client.Projects.CreateProject(projectOptions)
	if err != nil {
		log.Fatalf("Failed to create GitLab repository: %v", err)
	}

	return project, res, nil
}