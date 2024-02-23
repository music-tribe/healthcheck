package healthcheck

import (
	"net/http"
	"testing"
	"time"

	logger "dev.azure.com/MusicTribe/MT_CLOUD/mcloud-logger.git"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("when the service name is missing, we should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewClient("")
		})
	})

	t.Run("when we provide a service name, it should be set in the underlying client", func(t *testing.T) {
		expect := "someService"
		cli := NewClient(expect)

		var actual string
		if c, ok := cli.(*client); ok {
			actual = c.serviceName
		}

		assert.Equal(t, expect, actual)
	})

	t.Run("when we create a default Client, we should expect the failureStatusCode to be 503", func(t *testing.T) {
		expect := http.StatusServiceUnavailable
		cli := NewClient("svcName")

		var actual int
		if c, ok := cli.(*client); ok {
			actual = c.failureStatusCode
		}

		assert.Equal(t, expect, actual)
	})

	t.Run("when we set the failureStatusCode to 400 via functional option, we should expect the clients failureStatusCode to be 400", func(t *testing.T) {
		expect := http.StatusBadRequest
		cli := NewClient("svcName", func(co *ClientOptions) { co.FailureStatusCode = expect })

		var actual int
		if c, ok := cli.(*client); ok {
			actual = c.failureStatusCode
		}

		assert.Equal(t, expect, actual)
	})

	t.Run("when we create a default client, we should expect a default logger and not a nil value", func(t *testing.T) {
		cli := NewClient("svcName")

		var actual HealthCheckLogger
		if c, ok := cli.(*client); ok {
			actual = c.logger
		}

		assert.NotNil(t, actual)
	})

	t.Run("when we create a default client, we should expect a default timeout of 25 seconds", func(t *testing.T) {
		expect := time.Second * 10
		cli := NewClient("svcName", func(co *ClientOptions) { co.Timeout = expect })

		var actual time.Duration
		if c, ok := cli.(*client); ok {
			actual = c.timeout
		}

		assert.Equal(t, expect, actual)
	})

	t.Run("when we set a timeout of 10 seconds, we should expect the timeout to be set to 10 seconds", func(t *testing.T) {
		expect := time.Second * 25
		cli := NewClient("svcName")

		var actual time.Duration
		if c, ok := cli.(*client); ok {
			actual = c.timeout
		}

		assert.Equal(t, expect, actual)
	})

	t.Run("when we set our own logger via functional option, we should expect that logger to be on the client object", func(t *testing.T) {
		expect := logger.New("", false).StdLog("test", "logger")
		cli := NewClient("svcName", func(co *ClientOptions) { co.Logger = expect })

		var actual HealthCheckLogger
		if c, ok := cli.(*client); ok {
			actual = c.logger
		}

		assert.Equal(t, expect, actual)
	})
}
