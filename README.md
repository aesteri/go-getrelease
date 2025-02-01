
# FORKED 


# go-getrelease

go-getrelease is a library for Go (golang) for downloading release assets from source control hosting 
sites like github, and gitlab. This library only supports downloading files from public
urls and hence private repositories are not yet supported.

Library supports following clients:
* [x] Github
* [x] Gitlab

## Installation and Usage

Installation can be done with a normal `go get`:

```
$ go get github.com/dhillondeep/go-getrelease
```

## Basic Usage Instructions
* Client is created for the site you would like to download release assets from
* Call the getter functions using the client and provide asset name and destination
* That asset is downloaded

### Creating Client

```go
// Github
client := getrelease.NewGithubClient(nil)

// Gitlab
client := getrelease.NewBasicAuthGitlabClient(nil, "https://gitlab.com", "username", "password")
```

#### Authentication
go-getrelease library clients have different authentication methods based on which client is being used. Overview is
provided for each client.

##### Github
Github client does not directly handle authentication. Instead, when creating a new client, pass an http.Client that 
can handle authentication for you. This can be done as follows:

```go
import "golang.org/x/oauth2"

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client := getrelease.NewGithubClient(tc, "someOwner", "someRepo")
}
```

##### Gitlab
Gitlab client provides three ways of authentication: `OAuth`, `Private Token` and `Basic Authorization`. You can create
the client accordingly using one of New functions. Examples:

```go
import "golang.org/x/oauth2"

func main() {
	token := "----Some OAuth Token-----"
	client := getrelease.NewOAuthGitlabClient(nil, getrelease.GitlabDefaultBaseURL, token)
}
```

### Downloading Latest Release Asset
Downloading an asset from some release tag requires the location where the asset will be downloaded, name of asset, 
owner, and repo If the asset is an archived file, it is automatically `unarchieved` but, this can be turned off by
specifying `options`. Options are explained later on in another section. Client provided will provide information
on where to get the release asset from. This function also returns the name of the asset.

```go
if assetName, err := GetLatestAsset(client, "./download", "file.txt", "someOwner", "someRepo"); err != nil {
    panic(err);
} else {
	fmt.Println(assetName)
}
```

### Downloading Tag Release Asset
Downloading an asset from some release tag requires the location where the asset will be downloaded, name of asset, 
owner, repo, and tag name. If the asset is an archived file, it is automatically `unarchieved` but, this can be 
turned off by specifying `options`. Options are explained later on in another section. Client provided will provide 
information on where to get the release asset from. This function also returns the name of asset.

```go
if assetName, err := GetTagAsset(client, "./download", "file.txt", "someOwner", "someRepo", "v1.0.0"); err != nil {
    panic(err);
}  else {
  	fmt.Println(assetName)
}
```

### Asset Names
Asset names are the names of the assets that will be retrieved from the release. These names can be regex and this
regex will be matched against all the assets available. For examples if asset in latest release is `asset_0.1.0.txt`
but, you do not know the version so, you can name the asset as `asset_.*_.txt` and it will automatically be matched.

Example:
```go
if _, err := getrelease.GetTagAsset(client, "./wio", "wio_.*_linux_64bit.tar.gz",
    "v0.9.0", func(config *getrelease.Configuration) error {
        config.Checksum = "asset:wio_.*_checksums.txt"
        return nil
    }); err != nil {
        panic(err)
}
````


## Options
go-getrelease provide a way of customizing how files are downloaded and few other features like checksum verification.
When a call is made to `GetLatestAsset` or `GetTagAsset`, variadic number of `Options` can be provided where config
can be modified accordingly.

Example:

```go
if _, err := GetTagAsset(client, "./download", "file.txt", "v1.0.0", func(config *getrelease.Configuration) error {
	config.Checksum = "md5:46798b5cfca45c46a84b7419f8b74735"
	return nil
}); err != nil {
	panic(error)
}
```

Library comes with few options out of box that can be used:

* `WithProgress` - Allows for a user to track the progress of a download
* `WithChecksum` - Allows for a user to provide checksum for the file
* `WithArchive` - Allows for a user to customize unarchiving while working with archived files.
* `WithFileName` - Allows for a user to rename downloaded file.

## Progress
While downloading files, a listener can be registered that will listen to progress. Library calls `TrackProgress` method
and provides it with `src`, `currentSize`, `totalSize`, and `stream`. This information can be used to create a progress
tracker. Example of this can be see in [progress_tracker](examples/progress_tracker) example folder.

## Archive
go-getrelease will automatically unarchive files into a file or directory based on the extension of the file being requested.
This works for both file and directory downloads. 

While calling the function, you can set `config.Archive` option to a value specifying the format of achieve. If this is
not provided, it will automatically use the extension of the path to see if it appears archived. Unarchiving can be 
explicitly disabled by setting the value of this filed to `false`.

The following archive formats are supported:

* `tar.gz` and `tgz`
* `tar.bz2` and `tbz2`
* `tar.xz` and `txz`
* `zip`
* `gz`
* `bz2`
* `xz`

## Checksum
For file downloads, go-getrelease can automatically verify a checksum for you. Note that checksumming only works for
downloading files, not directories.

To checksum a file, set the value of `config.Checksum` option. The value must in the format of `type:value` or just
`value`, where type is `"md5"`, `"sha1"`, `"sha256"`, `"sha512"`, `"asset"`, `"link"`, `"file"`. The "value" must be
the actual checksum value of if type is "asset", "link", or "file", following actions are taken:

* "asset" - name of the asset inside the same release. It will be downloaded and checksum value will be used from that
* "link" - url link to some file and this file will be downloaded and checksum is used
* "file" - local file location from where checksum can be read
 
If no `type` is provided, error is thrown. Examples:

```go
if _, err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Checksum = "md5:46798b5cfca45c46a84b7419f8b74735"
	return nil
}); err != nil {
	panic(error)
}
```

```go
if _, err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Checksum = "asset:checksum.txt"
	return nil
}); err != nil {
	panic(error)
}
```

```go
if _, err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Checksum = "link:https://somelink/checksum.txt"
	return nil
}); err != nil {
	panic(error)
}
```

For checksum files, content of files are expected to be BSD or GNU style. Once go-getrelease is done with the checksum file;
it is deleted.

## FileName
You can rename the file being downloaded. Set `config.FileName` option value to new name.

## Pwd
When you are refering to local paths, the library needs to know the present working directory. This can be specified
by setting `config.Pwd` option value to a path. If this is set, you can use relative paths like for example:

```go
if _, err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Pwd = "someDir"
	config.Checksum = "file:./checksum.txt"
	return nil
}); err != nil {
	panic(error)
}
```

You can look at examples in [examples directory](examples). If you have any question, find a bug, create an issue. Feel free to
contribute to this project!
