package getrelease

type Configuration struct {
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

type Options func(config *Configuration) error

// configure configures provided options
func (config *Configuration) configure(opts ...Options) error {
	for _, opt := range opts {
		err := opt(config)
		if err != nil {
			return err
		}
	}
	return nil
}

// WithProgress allows for a user to track the progress of a download.
func WithProgress(pl ProgressTracker) func(*Configuration) error {
	return func(config *Configuration) error {
		config.ProgressTracker = pl
		return nil
	}
}

// WithChecksum allows for a user to provide checksum for the file.
func WithChecksum(checksum string) func(*Configuration) error {
	return func(config *Configuration) error {
		config.Checksum = checksum
		return nil
	}
}

// WithArchive allows for a user to customize unarchiving while working with archived files.
func WithArchive(archive string) func(*Configuration) error {
	return func(config *Configuration) error {
		config.Archive = archive
		return nil
	}
}

// WithFileName allows for a user to rename downloaded file.
func WithFileName(fileName string) func(*Configuration) error {
	return func(config *Configuration) error {
		config.Archive = fileName
		return nil
	}
}
