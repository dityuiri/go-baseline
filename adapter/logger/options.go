package logger

// Options represents options that can be used to configure the logger.
type Options struct {
	name   *string
	level  *LevelType
	output *OutputType

	skip *int // Number of stack frames to skip.

	isNoop bool
}

type Option func(options *Options)

func WithAppName(name string) Option {
	return func(options *Options) {
		options.name = &name
	}
}

func WithLevel(level LevelType) Option {
	return func(options *Options) {
		options.level = &level
	}
}

func WithOutput(output OutputType) Option {
	return func(options *Options) {
		options.output = &output
	}
}

func WithSkip(skip int) Option {
	return func(options *Options) {
		options.skip = &skip
	}
}

func WithNoOperation() Option {
	return func(options *Options) {
		options.isNoop = true
	}
}
