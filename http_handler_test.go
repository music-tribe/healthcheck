package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_client_HttpHandler(t *testing.T) {
	e := echo.New()
	t.Run("when the context is already cancelled, we should return that error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.TODO())
		cancel()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		cli := NewClient("tst")

		err := cli.HttpHandler(ctx, Test{name: "test", method: func(ctx context.Context) error { return nil }})(c)
		assert.ErrorContains(t, err, context.Canceled.Error())
		assert.Equal(t, 503, getStatusCode(rec, err))
	})

	t.Run("when a test errors, we should return that error", func(t *testing.T) {
		ctx := context.TODO()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		cli := NewClient("tst")
		expect := errors.New("some error")

		et := Test{
			name: "errors",
			method: func(ctx context.Context) error {
				return expect
			},
		}

		err := cli.HttpHandler(ctx, et)(c)
		assert.ErrorContains(t, err, expect.Error())
		assert.Equal(t, 503, getStatusCode(rec, err))
	})

	t.Run("when no tests error, we should return a 200", func(t *testing.T) {
		ctx := context.TODO()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		cli := NewClient("tst")

		err := cli.HttpHandler(
			ctx,
			Test{name: "t1", method: func(ctx context.Context) error { return nil }},
			Test{name: "t2", method: func(ctx context.Context) error { return nil }},
			Test{name: "t3", method: func(ctx context.Context) error { return nil }},
		)(c)

		assert.NoError(t, err)
		assert.Equal(t, 200, getStatusCode(rec, err))
	})
}

func getStatusCode(rec *httptest.ResponseRecorder, err error) int {
	ec := 500
	if err == nil {
		return rec.Code
	}

	httperr := &echo.HTTPError{}
	if errors.As(err, &httperr) {
		ec = httperr.Code
	}

	return ec
}
