package main

import (
	"github.com/dhillondeep/go-getrelease"
)

func main() {
	client := getrelease.NewGithubClient(nil)

	// progress bar extends the ProgressTracker
	if _, err := getrelease.GetTagAsset(client, "./download", "wio", "wio",
		"wio_.*_checksums.txt", "v0.9.0", getrelease.WithProgress(defaultProgressBar)); err != nil {
		panic(err)
	}
}
