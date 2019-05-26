package getrelease

import "fmt"

func gitlabError(err error) error {
	return fmt.Errorf("gitlab: %v", err)
}

func githubError(err error) error {
	return fmt.Errorf("github: %v", err)
}
