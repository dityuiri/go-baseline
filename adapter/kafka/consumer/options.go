package consumer

// Options represents options that can be used.
type Options struct {
	groupID     *string
	startOffset *int64
}

type Option func(options *Options)

func WithGroupID(groupID string) Option {
	return func(options *Options) {
		options.groupID = &groupID
	}
}
