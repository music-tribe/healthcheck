package healthcheck

import (
	"context"

	"github.com/labstack/echo/v4"
)

func (cli *client) HttpHandler(ctx context.Context, tests ...Test) func(echo.Context) error {
	return func(c echo.Context) error {
		err := laboratory(ctx, cli.logger, tests...)
		if err != nil {
			return echo.NewHTTPError(cli.failureStatusCode, err)
		}
		return err
	}
}
