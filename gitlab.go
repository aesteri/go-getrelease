package getrelease

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"path"
	"regexp"
)

const (
	GitlabDefaultBaseURL = "https://gitlab.com/"
)

// GitlabClient is a client for gitlab
type gitlabClient struct {
	httpClient *http.Client
	rawClient  *gitlab.Client
}

// NewOAuthGitlabClient provides connection with gitlab api through OAuth
func NewOAuthGitlabClient(client *http.Client, baseUrl, token string) (*gitlabClient, error) {
	rawClient := gitlab.NewOAuthClient(client, token)
	if err := rawClient.SetBaseURL(baseUrl); err != nil {
		return nil, gitlabError(err)
	}

	return &gitlabClient{
		httpClient: client,
		rawClient:  rawClient,
	}, nil
}

// NewBasicAuthGitlabClient provides connection with gitlab api through basic authentication
func NewBasicAuthGitlabClient(client *http.Client, baseUrl, username, password string) (*gitlabClient, error) {
	rawClient, err := gitlab.NewBasicAuthClient(client, baseUrl, username, password)
	if err != nil {
		return nil, gitlabError(err)
	}

	return &gitlabClient{
		httpClient: client,
		rawClient:  rawClient,
	}, nil
}

// NewPrivateTokenGitlabClient provides connection with gitlab api through private token
func NewPrivateTokenGitlabClient(client *http.Client, baseUrl, token string) (*gitlabClient, error) {
	rawClient := gitlab.NewClient(client, token)
	if err := rawClient.SetBaseURL(baseUrl); err != nil {
		return nil, gitlabError(err)
	}

	return &gitlabClient{
		httpClient: client,
		rawClient:  rawClient,
	}, nil
}

// createPid creates PID from owner and repo
func (client *gitlabClient) createPid(owner, repo string) string {
	return path.Clean(fmt.Sprintf("%s/%s", owner, repo))
}

// getAssetUrl is a generic url getter for github release assets
func (client *gitlabClient) getAssetUrl(release *gitlab.Release,
	response *gitlab.Response, assetNameReg *regexp.Regexp) (*http.Client, *string, error) {
	if response.StatusCode != 200 {
		return nil, nil, gitlabError(fmt.Errorf("invalid response status code: %s", response.Status))
	}

	for _, asset := range release.Assets.Links {
		if assetNameReg.Match([]byte(asset.Name)) {
			return client.httpClient, &asset.URL, nil
		}
	}

	return nil, nil, gitlabError(fmt.Errorf("%s release asset not found", assetNameReg.String()))
}

// getTagAssetUrl fetches gitlab releases of the project and returns link to specified asset of tag release
func (client *gitlabClient) getTagAssetUrl(assetNameReg *regexp.Regexp, owner, repo, tag string) (*http.Client, *string, error) {
	pid := client.createPid(owner, repo)
	release, response, err := client.rawClient.Releases.GetRelease(pid, tag)
	if response != nil && response.StatusCode == 403 {
		return nil, nil, gitlabError(fmt.Errorf("no release found for tag %s in %s", tag, pid))
	}
	if err != nil {
		return nil, nil, gitlabError(err)
	}

	return client.getAssetUrl(release, response, assetNameReg)
}

// getLatestAssetUrl fetches gitlab releases of the project and returns link to specified asset of latest release
func (client *gitlabClient) getLatestAssetUrl(assetNameReg *regexp.Regexp, owner, repo string) (*http.Client, *string, error) {
	pid := client.createPid(owner, repo)
	releases, response, err := client.rawClient.Releases.ListReleases(pid, nil)
	if err != nil {
		return nil, nil, gitlabError(err)
	}

	if len(releases) <= 0 {
		return nil, nil, gitlabError(fmt.Errorf("no release found for %s", pid))
	}

	return client.getAssetUrl(releases[0], response, assetNameReg)
}
