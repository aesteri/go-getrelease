package main

type Option struct {
	ProgressBar ProgressTracker
	Archive     string
	Checksum    string
	FileName    string
	Pwd         string
}

type Options func(option *Option) error

// WithProgress allows for a user to track the progress of a download.
func WithProgressBar(pl ProgressTracker) func(*Option) error {
	return func(option *Option) error {
		option.ProgressBar = pl
		return nil
	}
}

// Configure configures options
func (option *Option) configure(opts ...Options) error {
	for _, opt := range opts {
		err := opt(option)
		if err != nil {
			return err
		}
	}
	return nil
}
