package main

import (
	"github.com/hashicorp/go-getter"
	"io"
	"regexp"
)


type Client interface {
	getLatestAssetUrl(*regexp.Regexp) (*string, error)
	getTagAssetUrl(*regexp.Regexp, string) (*string, error)
}

// ProgressTracker allows to track the progress of downloads.
type ProgressTracker interface {
	TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser)
}

type urlGetter = func(assetNameReg *regexp.Regexp) (*string, error)

func get(dst string, assetNameReg *regexp.Regexp, urlGetter urlGetter, opts ...Options) error {
	option := &Option{}
	if err:= option.configure(opts...); err != nil {
		return err
	}

	retUrl, err := urlGetter(assetNameReg)
	if err != nil {
		return err
	}

	url, err := adjustUrlForGetter(*retUrl, urlGetter, option)
	if err != nil {
		return err
	}

	if err := getter.GetAny(dst, url, func(client *getter.Client) error {
		client.ProgressListener = option.ProgressBar
		client.Pwd = option.Pwd
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func GetTagAsset(client Client, dst, assetNameReg, tag string, opts ...Options) error {
	reg, err := regexp.Compile(assetNameReg)
	if err != nil {
		return err
	}

	urlGetter := func(assetNameReg *regexp.Regexp) (*string, error) {
		url, err := client.getTagAssetUrl(assetNameReg, tag)
		if err != nil {
			return nil, err
		}

		return url, nil
	}

	return get(dst, reg, urlGetter, opts...)
}

func GetLatestAsset(client Client, dst, assetNameReg string, opts ...Options) error {
	reg, err := regexp.Compile(assetNameReg)
	if err != nil {
		return err
	}

	urlGetter := func(assetNameReg *regexp.Regexp) (*string, error) {
		url, err := client.getLatestAssetUrl(assetNameReg)
		if err != nil {
			return nil, err
		}

		return url, nil
	}

	return get(dst, reg, urlGetter, opts...)
}
