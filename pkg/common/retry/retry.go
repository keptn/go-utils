package retry

import (
	"context"
	"fmt"
	"time"
)

const (
	DefaultNumberOfRetries     = 20
	DefaultDelayBetweenRetires = time.Second * 5
)

type Option func(*RetryConfiguration)

func NumberOfRetries(n uint) Option {
	return func(c *RetryConfiguration) {
		c.numberOfRetries = n
	}
}

func DelayBetweenRetries(d time.Duration) Option {
	return func(c *RetryConfiguration) {
		c.delayBetweenRetries = d
	}
}

func Context(ctx context.Context) Option {
	return func(c *RetryConfiguration) {
		c.context = ctx
	}
}

type RetryConfiguration struct {
	context             context.Context
	numberOfRetries     uint
	delayBetweenRetries time.Duration
}

type RetryFunc func() error

// Retry executes the retryFunc repeatedly until it was successful or canceled by the context
// The default number of retries is 20 and the default delay between retries is 5 seconds
func Retry(retryFunc RetryFunc, opts ...Option) error {
	configuration := &RetryConfiguration{
		numberOfRetries:     DefaultNumberOfRetries,
		delayBetweenRetries: DefaultDelayBetweenRetires,
		context:             context.TODO()}
	for _, opt := range opts {
		opt(configuration)
	}

	var i uint
	for i < configuration.numberOfRetries {
		err := retryFunc()
		if err != nil {
			select {
			case <-time.After(configuration.delayBetweenRetries):
			case <-configuration.context.Done():
				return fmt.Errorf("retry cancelled")
			}
		} else {
			return nil
		}
		i++
	}
	return fmt.Errorf("operation unsuccessful after %d retry", i)
}
