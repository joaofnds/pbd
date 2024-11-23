package config

import (
	"app/adapter/featureflags"
	"app/adapter/http"
	"app/adapter/metrics"
	"app/adapter/mongo"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/adapter/tracing"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("config", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewAppConfig),
		fx.Provide(func(config AppConfig) http.Config { return config.HTTP }),
		fx.Provide(func(config AppConfig) metrics.Config { return config.Metrics }),
		fx.Provide(func(config AppConfig) postgres.Config { return config.Postgres }),
		fx.Provide(func(config AppConfig) mongo.Config { return config.Mongo }),
		fx.Provide(func(config AppConfig) redis.Config { return config.Redis }),
		fx.Provide(func(config AppConfig) featureflags.Config { return config.FeatureFlags }),
		fx.Provide(func(config AppConfig) tracing.Config { return config.Tracing }),
	)
	Invokes = fx.Options(
		fx.Invoke(LoadConfig),
		fx.Invoke(ValidateConfig),
	)
)

type AppConfig struct {
	Env          string              `mapstructure:"env" validate:"required,oneof=development staging production"`
	HTTP         http.Config         `mapstructure:"http" validate:"required"`
	Metrics      metrics.Config      `mapstructure:"metrics" validate:"required"`
	Postgres     postgres.Config     `mapstructure:"postgres" validate:"required"`
	Mongo        mongo.Config        `mapstructure:"mongo" validate:"required"`
	Redis        redis.Config        `mapstructure:"redis" validate:"required"`
	FeatureFlags featureflags.Config `mapstructure:"feature_flags" validate:"required"`
	Tracing      tracing.Config      `mapstructure:"tracing" validate:"required"`
}

func LoadConfig() error {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		bindEnvs()
		return nil
	}

	viper.SetConfigFile(configFile)
	return viper.ReadInConfig()
}

func ValidateConfig(config AppConfig, validator *validator.Validate) error {
	return validator.Struct(config)
}

func NewAppConfig() (AppConfig, error) {
	var config AppConfig
	return config, viper.UnmarshalExact(&config)
}

func bindEnvs() {
	viper.MustBindEnv("env", "ENV")
	viper.MustBindEnv("metrics.address", "METRICS_ADDRESS")
	viper.MustBindEnv("http.port", "HTTP_PORT")
	viper.MustBindEnv("http.limiter.requests", "HTTP_LIMITER_REQUESTS")
	viper.MustBindEnv("http.limiter.expiration", "HTTP_LIMITER_EXPIRATION")
	viper.MustBindEnv("postgres.uri", "POSTGRES_URI")
	viper.MustBindEnv("mongo.uri", "MONGO_URI")
	viper.MustBindEnv("redis.addr", "REDIS_ADDR")
}
