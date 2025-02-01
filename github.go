package getrelease

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"regexp"
)

// GithubClient is a client for github
type GithubClient struct {
	httpClient *http.Client
	rawClient  *github.Client
}

// NewGithubClient returns github client
func NewGithubClient(client *http.Client) *GithubClient {
	return &GithubClient{
		httpClient: client,
		rawClient:  github.NewClient(client),
	}
}

// getAssetUrl is a generic url getter for github release assets
func (client *GithubClient) getAssetUrl(release *github.RepositoryRelease,
	response *github.Response, assetNameReg *regexp.Regexp) (*http.Client, *string, *string, error) {
	if response.StatusCode != 200 {
		return nil, nil, nil, githubError(fmt.Errorf("invalid response status code: %s", response.Status))
	}

	for _, asset := range release.Assets {
		if assetNameReg.Match([]byte(*asset.Name)) {
			return client.httpClient, asset.Name, asset.BrowserDownloadURL, nil
		}
	}

	return nil, nil, nil, githubError(fmt.Errorf("%s release asset not found", assetNameReg.String()))
}

// getTagAssetUrl fetches github releases of the project and returns link to specified asset of tag release
func (client *GithubClient) getTagAssetUrl(assetNameReg *regexp.Regexp,
	owner, repo, tag string) (*http.Client, *string, *string, error) {
	release, response, err := client.rawClient.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	if err != nil {
		return nil, nil, nil, githubError(err)
	}

	return client.getAssetUrl(release, response, assetNameReg)
}

// getLatestAssetUrl fetches github releases of the project and returns link to specified asset of latest release
func (client *GithubClient) getLatestAssetUrl(assetNameReg *regexp.Regexp,
	owner, repo string) (*http.Client, *string, *string, error) {
	release, response, err := client.rawClient.Repositories.GetLatestRelease(
		context.Background(), owner, repo)
	if err != nil {
		return nil, nil, nil, githubError(err)
	}

	return client.getAssetUrl(release, response, assetNameReg)
}
