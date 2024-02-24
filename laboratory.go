package healthcheck

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Test struct {
	name   string
	method TestMethod
}

type TestMethod func(context.Context) error

func NewTest(name string, method TestMethod) Test {
	return Test{name, method}
}

func laboratory(ctx context.Context, logger HealthCheckLogger, tests ...Test) error {
	if logger == nil {
		panic("laboratory: logger arg (type HealthCheckLogger) has nil value")
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, test := range tests {
		test := test
		g.Go(func() error {
			select {
			case <-ctx.Done():
				logger.Errorf("healthcheck Test '%s' not run due to %v", test.name, ctx.Err())
				return ctx.Err()
			default:
				err := test.method(ctx)
				if err != nil {
					logger.Errorf("healthcheck Test '%s' failed with error: %v", test.name, err)
				}
				return err
			}
		})
	}

	if err := g.Wait(); err != nil {
		logger.Errorf("laboratory test suite failed on: %v", err)
		return err
	}

	return nil
}
