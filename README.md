# healthcheck
A Golang health check client that exposes an echo handler for concurrently handling dependency health checks

## In Use

### Initializing the client and setting up an echo http handler
Here's a very basic example of how we might set up the healthcheck
```golang
    import github.com/music-tribe/healthcheck

    hcClient := healthcheck.NewClient("<myServiceName>")

    db := database.New()
    app := application.New(db)

    e := echo.New()

    appHealthCheckTest := healthcheck.NewTest(
        "app_health", 
        func(ctx context.Context) error {
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
                // the app has it's own simple health check function
                return app.Health()
        },
    )

    // imagine each dependency we want to test exposes a healthcheck test
    ctx := context.Background()
    e.GET("/readiness", hcClient.HttpHandler(ctx, appHealthCheckTest, db.HealthCheckTest)
```

### Creating a healthcheck test
Healthcheck tests are simple to create. They're just objects that contain a name and a simple method for checking
some parameter within your dependency. Here's a simple way to write a database connection test...
```golang
    type Database struct {...}

    func (db *Database) HealthCheckTest() healthcheck.Test {
        method := func(ctx context.Context) error {
            return db.session.Ping()
        }

        return healthcheck.NewTest("database ping", method)
    }

    ...

    // the test would be injected as below
    ctx := context.Background()
    e.GET("/readiness", hcClient.HttpHandler(ctx, db.HealthCheckTest())

```
