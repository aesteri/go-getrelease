# go-getrelease

go-getrelease is a library for Go (golang) for downloading release assets from source control hosting 
sites like github, gitlab, bitbucket, etc.

Library supports following clients:
* [x] Github
* [ ] Gitlab
* [ ] Bitbucket

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
client := &GithubClient{
    Owner: "someOwner",
    Repo: "someRepo",
}
```

Clients contains information required to fetch information from that source.

### Downloading Latest Release Asset
Downloading an asset from the latest release requires the name of asset and the location where the asset will be 
downloaded. If the asset is an archived file, it is automatically `unarchieved` but, this can be turned off by
specifying `options`. Options are explained later on in another section. Client provided will provide information
on where to get the release asset from.

```go
if err := GetLatestAsset(client, "./download", "file.txt"); err != nil {
    panic(err);
}
```

### Downloading Tag Release Asset
Downloading an asset from some release tag requires the tag name, name of asset, and the location where the asset will 
be downloaded. If the asset is an archived file, it is automatically `unarchieved` but, this can be turned off by
specifying `options`. Options are explained later on in another section. Client provided will provide information
on where to get the release asset from.

```go
if err := GetTagAsset(client, "./download", "v1.0.0", "file.txt"); err != nil {
    panic(err);
}
```

## Options
go-getrelease provide a way of customizing how files are downloaded and few other features like checksum verification.
When a call is made to `GetLatestAsset` or `GetTagAsset`, variadic number of `Options` can be provided where config
can be modified accordingly.

Example:

```go
if err := GetTagAsset(client, "./download", "v1.0.0", "file.txt", func(config *getrelease.Configuration) error {
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
tracker. Example of this can be see in [progress_tracker](examples/progress_tarcker) example folder.

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
if err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Checksum = "md5:46798b5cfca45c46a84b7419f8b74735"
	return nil
}); err != nil {
	panic(error)
}
```

```go
if err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Checksum = "asset:checksum.txt"
	return nil
}); err != nil {
	panic(error)
}
```

```go
if err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
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
if err := GetLatestAsset(client, "./download", "file.txt", func(config *getrelease.Configuration) error {
	config.Pwd = "someDir"
	config.Checksum = "file:./checksum.txt"
	return nil
}); err != nil {
	panic(error)
}
```

You can look at examples in [examples directory](examples). If you have any question, find a bug, create an issue. Feel free to
contribute to this project!
