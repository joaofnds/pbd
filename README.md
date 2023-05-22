# Features

- [Dependency Injection](cmd/app/app.go#L20) with [Fx](https://github.com/uber-go/fx)
- [Configuration](config/config.go#L47) with [Viper](https://github.com/spf13/viper)
- [Logging](adapters/logger/logger.go#L10) with [Zap](https://github.com/uber-go/zap)
- [Metrics](adapters/metrics/metrics.go#L22) with [Prometheus](https://github.com/prometheus/client_golang)
- [Health checks](adapters/health/controller.go#L18)
- [HTTP](adapters/http/fiber.go#L26) with [Fiber](https://github.com/gofiber/fiber)
- [Background](adapters/queue/client.go#L12) [tasks](user/queue/greeter.go#L27)/[workers](cmd/worker/worker.go#L14) with [Asynq](https://github.com/hibiken/asynq)
- [Testing](user/service_test.go#L70) with [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega)
- [Migrations](cmd/migrate/migrate.go#L20) with [Goose](https://github.com/pressly/goose)
- [Sto](user/mongo_repository.go)[ra](kv/redis_store.go)[ge](user/postgres_repository.go) with [Mongo](https://github.com/mongodb/mongo-go-driver), [Redis](https://github.com/redis/go-redis), and [Gorm](https://github.com/go-gorm/gorm) ([Postgres](https://github.com/go-gorm/postgres))
- [Version](.github/workflows/commit.yaml#L66) [management](.releaserc.yaml) with [Semantic Release](https://github.com/semantic-release/semantic-release)
- [Image](Dockerfile) (using [distroless](https://github.com/GoogleContainerTools/distroless)) [publishing to GitHub Container Registry](.github/workflows/build.yaml)
