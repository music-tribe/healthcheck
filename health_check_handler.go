package healthcheck

import "github.com/labstack/echo/v4"

type HealthChecker interface {
	HealthCheck() error
}

func (c *client) Handler(healthchecks ...HealthChecker) func(echo.Context) error {
	return func(ctx echo.Context) error {
		return nil
	}
}
