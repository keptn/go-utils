package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfigurableSleeper(t *testing.T) {
	timeToSleep := 5 * time.Second

	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{timeToSleep, spyTime.Sleep}
	sleeper.Sleep()
	assert.Equal(t, timeToSleep, spyTime.durationSlept)
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}
