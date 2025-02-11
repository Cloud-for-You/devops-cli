package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Member struct {
	Name string
}

// currentMembers is members in GitLab group
// desiredMembers is members in SRC group (LDAP, Azure, File, etc...)
// missing: User is LDAP but not in GitLab (Create)
// extra: User in GitLab but not in LDAP (Delete)
func CompareMembers(gitlabMembers, sourceMembers []Member) (missing, extra []Member) {
	// Vytvorime mapy pro rychle vyhledavani
	gitlabSet := make(map[string]struct{})
	srcSet := make(map[string]struct{})

	// Naplnime mapu
	for _, m := range gitlabMembers {
		gitlabSet[m.Name] = struct{}{}
	}
	for _, m := range sourceMembers {
		srcSet[m.Name] = struct{}{}
	}

	// Najdeme chybejici cleny (v source, ale ne v GitLab)
	for _, m := range sourceMembers {
		if _, exists := gitlabSet[m.Name]; !exists {
			missing = append(missing, m)
		}
	}
	fmt.Printf("Created: %v\n", missing)

	// Najdeme prebyvajici cleny (v GitLab, ale ne v source)
	for _, m := range gitlabMembers {
		if _, exists := srcSet[m.Name]; !exists {
			extra = append(extra, m)
		}
	}
	fmt.Printf("Deleted: %v\n", extra)

	return missing, extra
}

func LoadConfigFile(path string) (map[string]interface{}, error) {
	var data map[string]interface{}

	var content []byte
	var err error

  if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// Načtení souboru z URL
		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch config from URL: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch config: HTTP %d", resp.StatusCode)
		}

		content, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
	} else {
		// Načtení lokálního souboru
		content, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %v", err)
		}
	}

	// Rozlišení formátu podle obsahu
	if json.Valid(content) {
		err = json.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}
	} else {
		err = yaml.Unmarshal(content, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %v", err)
		}
	}
	
	return data, nil
}

func Ptr(s string) *string {
	return &s
}
