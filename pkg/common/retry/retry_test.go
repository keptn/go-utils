package retry_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDefaultNumberOfRetriesAllFail(t *testing.T) {
	var count int
	err := retry.Retry(
		func() error {
			count++
			fmt.Println(count)
			return errors.New("test")
		},
		retry.DelayBetweenRetries(time.Millisecond*10),
	)
	assert.Equal(t, retry.DefaultNumberOfRetries, count)
	assert.NotNil(t, err)
}

func TestDefaultLastRetrySucceeded(t *testing.T) {
	var count int
	err := retry.Retry(
		func() error {
			count++
			if count == retry.DefaultNumberOfRetries {
				return nil
			}
			return errors.New("test")
		},
		retry.DelayBetweenRetries(time.Millisecond*10),
	)
	assert.Equal(t, retry.DefaultNumberOfRetries, count)
	assert.Nil(t, err)
}

func TestCustomNumberOfRetriesAllFail(t *testing.T) {
	var count int
	err := retry.Retry(
		func() error {
			count++
			fmt.Println(count)
			return errors.New("test")
		},
		retry.DelayBetweenRetries(time.Millisecond*1),
		retry.NumberOfRetries(50),
	)
	assert.Equal(t, 50, count)
	assert.NotNil(t, err)
}

func TestFirstCallSucceeded(t *testing.T) {
	var count int
	err := retry.Retry(
		func() error {
			count++
			return nil
		},
		retry.DelayBetweenRetries(time.Millisecond*10),
	)
	assert.Equal(t, 1, count)
	assert.Nil(t, err)
}

func TestCancelRetries(t *testing.T) {
	var count int
	ctx, cancel := context.WithCancel(context.TODO())
	err := retry.Retry(
		func() error {
			count++
			if count > 10 {
				cancel()
			}
			return errors.New("test")
		},
		retry.DelayBetweenRetries(time.Millisecond*10),
		retry.Context(ctx),
	)
	assert.Equal(t, 11, count)
	assert.NotNil(t, err)
}
