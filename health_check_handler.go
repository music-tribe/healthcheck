package healthcheck

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
)

type HealthChecker interface {
	HealthCheck(ctx context.Context, logger HealthCheckLogger) <-chan error
}

func (cli *client) Handler(ctx context.Context, healthchecks ...HealthChecker) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := context.WithDeadline(ctx, time.Now().Add(cli.timeout))
		defer cancel()
		g, ctx := errgroup.WithContext(ctx)

		for _, check := range healthchecks {
			// shadow check inside the loop before attempting to reference it inside the goroutine
			_check := check
			g.Go(func() error {
				select {
				case err, ok := <-_check.HealthCheck(ctx, cli.logger):
					if !ok {
						cli.logger.Infof("healthcheck.Handler: Healthcheck err channel not OK - may have been closed already")
						return nil
					}
					return err
				case <-ctx.Done():
					return ctx.Err()
				}
			})
		}

		if err := g.Wait(); err != nil {
			cli.logger.Errorf("healthcheck.Handler: %v", err)
			return err
		}

		return nil
	}
}
