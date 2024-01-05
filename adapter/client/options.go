package client

type Option func(*Client)

func WithDoer(doer Doer) Option {
	return func(c *Client) {
		c.doer = doer
	}
}

func WithMaxRetry(maxRetry int) Option {
	return func(c *Client) {
		c.maxRetry = maxRetry
	}
}
