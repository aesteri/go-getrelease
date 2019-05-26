package getrelease

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"regexp"
)

var rawClient = github.NewClient(nil)

// GithubClient is a client for github
type GithubClient struct {
	Owner string
	Repo  string
}

// getAssetUrl is a generic url getter for github release assets
func getAssetUrl(release *github.RepositoryRelease, response *github.Response, assetNameReg *regexp.Regexp) (*string, error) {
	if response.StatusCode != 200 {
		return nil, errors.New("invalid response status code: " + response.Status)
	}

	for _, asset := range release.Assets {
		if assetNameReg.Match([]byte(*asset.Name)) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("%s release asset not found", assetNameReg.String()))
}

// getTagAssetUrl fetches github releases of the project and returns link to specified asset of tag release
func (client *GithubClient) getTagAssetUrl(assetNameReg *regexp.Regexp, tag string) (*string, error) {
	release, response, err := rawClient.Repositories.GetReleaseByTag(
		context.Background(), client.Owner, client.Repo, tag)
	if err != nil {
		return nil, err
	}

	return getAssetUrl(release, response, assetNameReg)
}

// getLatestAssetUrl fetches github releases of the project and returns link to specified asset of latest release
func (client *GithubClient) getLatestAssetUrl(assetNameReg *regexp.Regexp) (*string, error) {
	release, response, err := rawClient.Repositories.GetLatestRelease(context.Background(), client.Owner, client.Repo)
	if err != nil {
		return nil, err
	}

	return getAssetUrl(release, response, assetNameReg)
}
