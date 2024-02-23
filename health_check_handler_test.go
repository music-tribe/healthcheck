package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_client_Handler(t *testing.T) {
	t.Run("when the client timeout is exceeded, the handler should pass this on via context to any running HealthCheckers", func(t *testing.T) {
		cli := NewClient("someService", func(co *ClientOptions) { co.Timeout = time.Second })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		err := cli.Handler(context.TODO(), &stubHealthChecker{delay: time.Second * 2})(c)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("when the a Healthchecker returns an error, the other Healthcheckers should get cancelled", func(t *testing.T) {
		_err := errors.New("some error")
		cli := NewClient("someService")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		hc1 := &stubHealthChecker{delay: time.Second * 1, err: _err}
		hc2 := &stubHealthChecker{delay: time.Second * 10}
		hc3 := &stubHealthChecker{delay: time.Second * 10}
		err := cli.Handler(context.TODO(), hc1, hc2, hc3)(c)

		assert.ErrorIs(t, err, _err)
	})

	t.Run("when none of the Healthcheckers return an error, we should receive no error", func(t *testing.T) {
		cli := NewClient("someService")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		hc1 := &stubHealthChecker{delay: time.Second * 1, err: nil}
		hc2 := &stubHealthChecker{delay: time.Second * 2}
		hc3 := &stubHealthChecker{delay: time.Second * 3}
		err := cli.Handler(context.TODO(), hc1, hc2, hc3)(c)

		assert.NoError(t, err)
	})
}

type stubHealthChecker struct {
	delay time.Duration
	err   error
}

func (shc *stubHealthChecker) HealthCheck(ctx context.Context, logger HealthCheckLogger) <-chan error {
	errChan := make(chan error)
	fnc := func() error {
		time.Sleep(shc.delay)
		return shc.err
	}

	go func() {
		defer close(errChan)
		select {
		case errChan <- fnc():
		default:
			// if the error channel is blocked, allow the goroutine to close
		}
	}()

	return errChan
}
