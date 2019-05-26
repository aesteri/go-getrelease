package main

import (
	"github.com/dhillondeep/go-getrelease"
)

func main() {
	// Basic Auth
	clientBasic, err := getrelease.NewBasicAuthGitlabClient(
		nil, getrelease.GitlabDefaultBaseURL, "username", "password")
	if err != nil {
		panic(err)
	}

	if err := getrelease.GetLatestAsset(clientBasic, "./download", "assetName",
		"owner", "repo"); err != nil {
		panic(err)
	}

	// OAuth
	clientOAuth, err := getrelease.NewOAuthGitlabClient(
		nil, getrelease.GitlabDefaultBaseURL, "someToken")
	if err != nil {
		panic(err)
	}

	if err := getrelease.GetLatestAsset(clientOAuth, "./download", "assetName",
		"owner", "repo"); err != nil {
		panic(err)
	}

	// Private Token
	clientPrivate, err := getrelease.NewPrivateTokenGitlabClient(
		nil, getrelease.GitlabDefaultBaseURL, "someToken")
	if err != nil {
		panic(err)
	}

	if err := getrelease.GetLatestAsset(clientPrivate, "./download", "assetName",
		"owner", "repo"); err != nil {
		panic(err)
	}
}
