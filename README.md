# healthcheck
A Golang health check client that exposes an echo handler for concurrently handling dependency health checks

## In Use

### Initializing the client
```golang
    import github.com/music-tribe/healthcheck

    hcClient := healthcheck.NewClient("<myServiceName>")

    e := echo.New()

    // for this example service, database and storage objects all satisfy HealthChecker interface

    e.GET("/readiness", hcClient.Handler(service, database, storage))
```

### Implementing the Healthchecker on one of our dependencies
```golang
    type Database struct {...}

    func (db *Database) HealthCheck(ctx context.Context, logger HealthCheckLogger) <-chan error {
        errChan := make(chan error)

        go func() {
            defer close(errChan)

            select {
            case errChan<-db.Ping():
            default:
                // use default incase err chan is blocked
            }
        }()

        return errChan
    }
