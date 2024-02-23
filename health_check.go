package healthcheck

import (
	"net/http"

	logger "dev.azure.com/MusicTribe/MT_CLOUD/mcloud-logger.git"
	"github.com/labstack/echo/v4"
)

type Client interface {
	Handler(healthchecks ...HealthChecker) func(echo.Context) error
}

type HealthCheckLogger interface {
	Infof(format string, items ...interface{})
	Errorf(format string, items ...interface{})
}

type client struct {
	failureStatusCode int
	serviceName       string
	logger            HealthCheckLogger
}

type ClientOptions struct {
	FailureStatusCode int
	Logger            HealthCheckLogger
}

type ClientOption func(*ClientOptions)

func NewClient(serviceName string, options ...ClientOption) Client {
	if serviceName == "" {
		panic("NewClient: missing serviceName")
	}

	ops := ClientOptions{
		FailureStatusCode: http.StatusServiceUnavailable,
		Logger:            logger.New("", false).StdLog(serviceName, "healthcheck.Handler"),
	}

	for _, optFunc := range options {
		optFunc(&ops)
	}

	return &client{
		failureStatusCode: ops.FailureStatusCode,
		serviceName:       serviceName,
		logger:            ops.Logger,
	}
}
