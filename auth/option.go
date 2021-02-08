package auth

import "time"

type Options func(*Option)

type Option struct {
	tokenTTL  time.Duration
	codeTTL   time.Duration
	sender    Sender
	render    Render
	filter    Filter
	generater Generater
}

func WithSender(sender Sender) Options {
	return func(o *Option) {
		o.sender = sender
	}
}

func WithRender(render Render) Options {
	return func(o *Option) {
		o.render = render
	}
}

func WithFilter(filter Filter) Options {
	return func(o *Option) {
		o.filter = filter
	}
}

func WithGenerater(generater Generater) Options {
	return func(o *Option) {
		o.generater = generater
	}
}

func WithCodeTTL(d time.Duration) Options {
	return func(o *Option) {
		o.codeTTL = d
	}
}

func WithTokenTTL(d time.Duration) Options {
	return func(o *Option) {
		o.tokenTTL = d
	}
}
