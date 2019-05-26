package getrelease

import (
	"github.com/hashicorp/go-getter"
	"io"
	"net/http"
	"regexp"
)

// Client is one of github, gitlab connection that allows for fetching assets from releases
type Client interface {
	getLatestAssetUrl(*regexp.Regexp, string, string) (*http.Client, *string, error)
	getTagAssetUrl(*regexp.Regexp, string, string, string) (*http.Client, *string, error)
}

// ProgressTracker allows to track the progress of downloads.
type ProgressTracker interface {
	TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser)
}

// urlGetter is a type that returns url fetched for the asset
type urlGetter = func(assetNameReg *regexp.Regexp) (*http.Client, *string, error)

// get is a generic asset fetcher and downloader
func get(dst string, assetNameReg *regexp.Regexp, urlGetter urlGetter, opts ...Options) error {
	config := &Configuration{}
	if err := config.configure(opts...); err != nil {
		return err
	}

	httpClient, retUrl, err := urlGetter(assetNameReg)
	if err != nil {
		return err
	}

	url, err := adjustUrlForGetter(*retUrl, urlGetter, config)
	if err != nil {
		return err
	}

	if err := getter.GetAny(dst, url, func(client *getter.Client) error {
		httpGetter := &getter.HttpGetter{
			Netrc:  true,
			Client: httpClient,
		}
		client.Getters["http"] = httpGetter
		client.Getters["https"] = httpGetter

		client.ProgressListener = config.ProgressTracker
		client.Pwd = config.Pwd
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// GetTagAsset fetches provided tag release and matches the provided asset name regex with all the assets.
// If asset is found, it is downloaded in provided dst otherwise error is thrown
func GetTagAsset(client Client, dst, assetNameReg, owner, repo, tag string, opts ...Options) error {
	reg, err := regexp.Compile(assetNameReg)
	if err != nil {
		return err
	}

	urlGetter := func(assetNameReg *regexp.Regexp) (*http.Client, *string, error) {
		client, url, err := client.getTagAssetUrl(assetNameReg, owner, repo, tag)
		if err != nil {
			return nil, nil, err
		}

		return client, url, nil
	}

	return get(dst, reg, urlGetter, opts...)
}

// GetLatestAsset fetches latest release and matches the provided asset name regex with all the assets.
// If asset is found, it is downloaded in provided dst otherwise error is thrown
func GetLatestAsset(client Client, dst, assetNameReg, owner, repo string, opts ...Options) error {
	reg, err := regexp.Compile(assetNameReg)
	if err != nil {
		return err
	}

	urlGetter := func(assetNameReg *regexp.Regexp) (*http.Client, *string, error) {
		client, url, err := client.getLatestAssetUrl(assetNameReg, owner, repo)
		if err != nil {
			return nil, nil, err
		}

		return client, url, nil
	}

	return get(dst, reg, urlGetter, opts...)
}
