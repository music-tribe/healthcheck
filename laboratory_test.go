package healthcheck

import (
	"context"
	"errors"
	"testing"
	"time"

	logger "dev.azure.com/MusicTribe/MT_CLOUD/mcloud-logger.git"
	"github.com/stretchr/testify/assert"
)

func Test_laboratory(t *testing.T) {
	tests := []Test{
		{name: "test1", method: func(ctx context.Context) error { time.Sleep(time.Second); return nil }},
		{name: "test2", method: func(ctx context.Context) error { time.Sleep(time.Second); return nil }},
		{name: "test3", method: func(ctx context.Context) error { time.Sleep(time.Second); return nil }},
	}

	logger := logger.New("", false).StdLog("test", "testing")

	t.Run("when the logger arg has a nil value we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			laboratory(context.TODO(), nil)
		})
	})

	t.Run("when the context is cancelled before our test can run, we should receive the context canceled error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()

		err := laboratory(ctx, logger, tests...)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("when the one of our test methods returns an error, we should receive that error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		expect := errors.New("some error")

		test := Test{
			name: "erroring",
			method: func(ctx context.Context) error {
				time.Sleep(200 * time.Millisecond)
				return expect
			},
		}

		actual := laboratory(ctx, logger, test)
		assert.ErrorIs(t, actual, expect)
	})

	t.Run("when the one of our test methods context is cancelled, we should receive that error", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(time.Second))
		defer cancel()

		test := Test{
			name: "erroring",
			method: func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 2):
					return nil
				}
			},
		}

		actual := laboratory(ctx, logger, test)
		assert.ErrorIs(t, actual, context.DeadlineExceeded)
	})

	t.Run("when none of the tests fail, we should receive no error", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(time.Second*5))
		defer cancel()

		actual := laboratory(ctx, logger, tests...)
		assert.NoError(t, actual)
	})
}
