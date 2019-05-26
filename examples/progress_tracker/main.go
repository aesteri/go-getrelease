package main

import (
	"os"
)

func main()  {
	client := &GithubClient{
		Owner: "wio",
		Repo: "wio",
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := GetTagAsset(client,
		"./download", "wio_.*_macOS_64bit.tar.gz", "v0.9.0", WithProgressBar(defaultProgressBar), func(option *Option) error {
			option.Archive = "false"
			option.FileName = "wio.tar.gz"
			option.Checksum = "link::https://github.com/wio/wio/releases/download/v0.9.0/wio_0.9.0_checksums.txt"
			option.Pwd = dir
			return nil
		}); err != nil {
		panic(err)
	}
}
