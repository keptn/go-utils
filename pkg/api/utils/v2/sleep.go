package v2

import "time"

// Sleeper defines the interface to sleep
type Sleeper interface {
	Sleep()
}

// ConfigurableSleeper is an implementation of a sleeper
// that can be configured to sleep for a specific duration
type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

// Sleep pauses the execution
func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

// NewConfigurableSleeper creates a new instance of a configurable sleeper which will pause execution
// of the current thread for a given duration
func NewConfigurableSleeper(duration time.Duration) *ConfigurableSleeper {
	return &ConfigurableSleeper{
		duration: duration,
		sleep:    time.Sleep,
	}
}

// FakeSleeper is a sleeper that does not sleep
type FakeSleeper struct {
}

// Sleep does nothing, not even sleep
func (f *FakeSleeper) Sleep() {
	// no-op
}

// NewFakeSleeper creates a new instance of a FakeSleeper
func NewFakeSleeper() *FakeSleeper {
	return &FakeSleeper{}
}
