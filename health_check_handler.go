package healthcheck

import (
	"context"

	"github.com/labstack/echo/v4"
)

type HealthChecker interface {
	HealthCheck() error
}

func (c *client) Handler(ctx context.Context, healthchecks ...HealthChecker) func(echo.Context) error {
	return func(ctx echo.Context) error {
		return nil
	}
}
