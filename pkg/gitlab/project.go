package gitlab

import (
	"fmt"
	"log"
	"net/http"

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

	project, res, err := client.Projects.CreateProject(projectOptions)
	if err != nil {
		log.Fatalf("Failed to create GitLab repository: %v", err)
	}

	// Create and Invite Groups for members roles
	// Maintainer
	if maintainerGroupName != nil && *maintainerGroupName != "" {
		maintainerGroup, res, err := CreateGroup(client, *maintainerGroupName, "", "private")
		if err != nil {
			if res != nil && res.StatusCode == http.StatusConflict {
				fmt.Printf("Group '%s' is exists.\n", *maintainerGroupName)
			} else {
				log.Fatalf("Failed to create GitLab group '%s': %v\n", *maintainerGroupName, err)
			}
		}

		err = InviteGroupToProject(client, project.ID, maintainerGroup.ID, gitlab.Ptr(gitlab.MaintainerPermissions))
		if err != nil {
			log.Fatalf("Failed to invite maintainer group to project: %v", err)
		}
	}

	// Developer
	if developerGroupName != nil && *developerGroupName != "" {
		developerGroup, res, err := CreateGroup(client, *developerGroupName, "", "private")
		if err != nil {
			if res != nil && res.StatusCode == http.StatusConflict {
				fmt.Printf("Group '%s' is exists.\n", *developerGroupName)
			} else {
				log.Fatalf("Failed to create GitLab group '%s': %v\n", *developerGroupName, err)
			}
		}

		err = InviteGroupToProject(client, project.ID, developerGroup.ID, gitlab.Ptr(gitlab.DeveloperPermissions))
		if err != nil {
			log.Fatalf("Failed to invite developer group to project: %v", err)
		}
	}

	return project, res, nil
}
