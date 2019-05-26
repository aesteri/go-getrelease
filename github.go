package getrelease

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"regexp"
)

// GithubClient is a client for github
type githubClient struct {
	owner      string
	repo       string
	httpClient *http.Client
	rawClient  *github.Client
}

// NewGithubClient returns github client
func NewGithubClient(client *http.Client, owner, repo string) *githubClient {
	return &githubClient{
		owner:      owner,
		repo:       repo,
		httpClient: client,
		rawClient:  github.NewClient(client),
	}
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
func (client *githubClient) getTagAssetUrl(assetNameReg *regexp.Regexp, tag string) (*string, error) {
	release, response, err := client.rawClient.Repositories.GetReleaseByTag(
		context.Background(), client.owner, client.repo, tag)
	if err != nil {
		return nil, err
	}

	return getAssetUrl(release, response, assetNameReg)
}

// getLatestAssetUrl fetches github releases of the project and returns link to specified asset of latest release
func (client *githubClient) getLatestAssetUrl(assetNameReg *regexp.Regexp) (*string, error) {
	release, response, err := client.rawClient.Repositories.GetLatestRelease(
		context.Background(), client.owner, client.repo)
	if err != nil {
		return nil, err
	}

	return getAssetUrl(release, response, assetNameReg)
}
