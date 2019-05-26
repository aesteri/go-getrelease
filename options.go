package getrelease

type Option struct {
	// ProgressTracker allows to track the progress of downloads.
	ProgressTracker ProgressTracker

	// Archive is used for specifying achieve type while unarchiving and "false" can be used to turn unarchiving off
	Archive     string

	// Checksum to verify the downloaded file
	Checksum    string

	// Filename to rename the downloaded file to
	FileName    string

	// Pwd is the working directory for detection. If this isn't set, some
	// detection may fail. Client will not default pwd to the current
	// working directory for security reasons.
	Pwd         string
}

type Options func(option *Option) error

// configure configures provided options
func (option *Option) configure(opts ...Options) error {
	for _, opt := range opts {
		err := opt(option)
		if err != nil {
			return err
		}
	}
	return nil
}

// WithProgress allows for a user to track the progress of a download.
func WithProgress(pl ProgressTracker) func(*Option) error {
	return func(option *Option) error {
		option.ProgressTracker = pl
		return nil
	}
}

// WithChecksum allows for a user to provide checksum for the file
func WithChecksum(checksum string) func(*Option) error {
	return func(option *Option) error {
		option.Checksum = checksum
		return nil
	}
}

// WithArchive allows for a user to customize unarchiving while working with archived files
func WithArchive(archive string) func(*Option) error {
	return func(option *Option) error {
		option.Archive = archive
		return nil
	}
}

// WithFileName allows for a user to rename downloaded file
func WithFileName(fileName string) func(*Option) error {
	return func(option *Option) error {
		option.Archive = fileName
		return nil
	}
}
