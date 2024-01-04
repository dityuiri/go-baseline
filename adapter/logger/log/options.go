package log

import (
	"fmt"
)

type Options struct {
	data *string
	err  *error

	skip *int // Number of stack frames to skip.
}

type Option func(options *Options)

func WithData(format string, args ...interface{}) Option {
	return func(options *Options) {
		data := fmt.Sprintf(format, args...)
		options.data = &data
	}
}

func WithError(err error) Option {
	return func(options *Options) {
		options.err = &err
	}
}

func WithSkip(skip int) Option {
	return func(options *Options) {
		options.skip = &skip
	}
}

func (options *Options) GetData() *string {
	return options.data
}

func (options *Options) GetError() *error {
	return options.err
}

func (options *Options) GetSkip() *int {
	return options.skip
}
