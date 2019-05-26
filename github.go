package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/github"
	"regexp"
)

var rawClient = github.NewClient(nil)

type GithubClient struct {
	Owner string
	Repo  string
}

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

func (client *GithubClient) getTagAssetUrl(assetNameReg *regexp.Regexp, tag string) (*string, error) {
	release, response, err := rawClient.Repositories.GetReleaseByTag(
		context.Background(), client.Owner, client.Repo, tag)
	if err != nil {
		return nil, err
	}

	return getAssetUrl(release, response, assetNameReg)
}

func (client *GithubClient) getLatestAssetUrl(assetNameReg *regexp.Regexp) (*string, error) {
	release, response, err := rawClient.Repositories.GetLatestRelease(context.Background(), client.Owner, client.Repo)
	if err != nil {
		return nil, err
	}

	return getAssetUrl(release, response, assetNameReg)
}
