package main

import (
	"github.com/dhillondeep/go-getrelease"
)

func main()  {
	client := getrelease.NewGithubClient(nil)

	// wio_0.9.0_linux_64bit.tar.gz file will be downloaded inside wio folder and later renamed to wio.tar.gz
	// The file won't be unarchieved and checksum will be verified
	if err := getrelease.GetTagAsset(client, "./wio", "wio", "wio", "wio_.*_linux_64bit.tar.gz",
		"v0.9.0", func(config *getrelease.Configuration) error {
			config.Checksum = "asset:wio_.*_checksums.txt"
			config.Archive = "false"
			config.FileName = "wio.tar.gz"
			return nil
		}, getrelease.WithProgress(defaultProgressBar)); err != nil {
			panic(err)
	}
}
