package main

import (
	"github.com/dhillondeep/go-getrelease"
)

func main() {
	client := getrelease.NewGithubClient(nil, "wio", "wio")

	// progress bar extends the ProgressTracker
	if err := getrelease.GetTagAsset(client, "./download",
		"wio_.*_checksums.txt", "v0.9.0", getrelease.WithProgress(defaultProgressBar)); err != nil {
		panic(err)
	}
}
