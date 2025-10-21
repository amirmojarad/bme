package conf

import (
	"bme/pkg/errorext"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

type AppEnv string

const (
	AppEnvDebug AppEnv = "DEBUG"
	AppEnvProd  AppEnv = "PROD"
)

type AppConfig struct {
	App struct {
		Env  AppEnv
		Port int
	}

	Jwt struct {
		AccessSecret         []byte
		RefreshSecret        []byte
		RefreshTokenDuration time.Duration
		AccessTokenDuration  time.Duration
		Issuer               string
	}

	Postgres struct {
		Username              string
		Password              string
		Host                  string
		Name                  string
		Port                  uint64
		SslMode               string
		Timezone              string
		MaxIdleConnections    int
		MaxOpenConnections    int
		ConnectionMaxLifetime time.Duration
		MigrationPath         string
	}
}

type TelegramBot struct {
	ApiKey      string
	BotUsername string
}

func New() (*AppConfig, error) {
	cfg := &AppConfig{}

	setJwt(cfg)

	if err := setApp(cfg); err != nil {
		return nil, err
	}

	if err := setDatabase(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setJwt(cfg *AppConfig) {
	cfg.Jwt.AccessSecret = []byte(os.Getenv("JWT_ACCESS_SECRET"))
	cfg.Jwt.RefreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
	cfg.Jwt.Issuer = os.Getenv("JWT_ISSUER")

	accessDuration, err := envConverter("ACCESS_TOKEN_DURATION", time.ParseDuration)
	if err != nil {
		cfg.Jwt.AccessTokenDuration = 30 * time.Minute
	} else {
		cfg.Jwt.AccessTokenDuration = accessDuration
	}

	refreshDuration, err := envConverter("REFRESH_TOKEN_DURATION", time.ParseDuration)
	if err != nil {
		cfg.Jwt.RefreshTokenDuration = 3 * 24 * time.Hour
	} else {
		cfg.Jwt.RefreshTokenDuration = refreshDuration
	}
}

func setApp(cfg *AppConfig) error {
	appEnv := os.Getenv("APP_ENV")

	cfg.App.Env = newAppEnv(appEnv)

	port, err := envConverter("APP_PORT", strconv.Atoi)
	if err != nil {
		return err
	}

	cfg.App.Port = port

	return nil
}

func setDatabase(cfg *AppConfig) error {
	cfg.Postgres.Username = os.Getenv("POSTGRES_DATABASE_USERNAME")
	cfg.Postgres.Password = os.Getenv("POSTGRES_DATABASE_PASSWORD")
	cfg.Postgres.Host = os.Getenv("POSTGRES_DATABASE_HOST")
	cfg.Postgres.Name = os.Getenv("POSTGRES_DATABASE_NAME")

	port, err := envConverter("POSTGRES_DATABASE_PORT", func(v string) (uint64, error) {
		return strconv.ParseUint(v, 10, 32)
	})
	if err != nil {
		return err
	}

	cfg.Postgres.Port = port

	cfg.Postgres.SslMode = os.Getenv("POSTGRES_DATABASE_SSLMODE")
	cfg.Postgres.Timezone = os.Getenv("POSTGRES_DATABASE_TIMEZONE")

	maxConn, err := envConverter("POSTGRES_DATABASE_MAX_OPEN_CONN", strconv.Atoi)
	if err != nil {
		return err
	}

	cfg.Postgres.MaxOpenConnections = maxConn

	maxIdle, err := envConverter("POSTGRES_DATABASE_MAX_IDLE_CONN", strconv.Atoi)
	if err != nil {
		return err
	}

	cfg.Postgres.MaxIdleConnections = maxIdle

	connMaxLif, err := envConverter("POSTGRES_DATABASE_CONN_MAX_LIFETIME", time.ParseDuration)
	if err != nil {
		return err
	}

	cfg.Postgres.ConnectionMaxLifetime = connMaxLif

	cfg.Postgres.MigrationPath = os.Getenv("POSTGRES_DATABASE_MIGRATION_PATH")

	return nil
}

func envConverter[T any](key string, converter func(v string) (T, error)) (T, error) {
	var noop T

	value := os.Getenv(key)

	if value == "" {
		return noop, errorext.NewValidation(errors.New(errorext.EnvValueIsEmpty.String()), errorext.EnvValueIsEmpty)
	}

	converted, err := converter(value)
	if err != nil {
		return noop, errorext.NewValidation(errors.New(errorext.EnvValueIsEmpty.String()), errorext.EnvValueIsEmpty)
	}

	return converted, nil
}

func newAppEnv(env string) AppEnv {
	switch env {
	case AppEnvDebug.String():
		return AppEnvDebug
	case AppEnvProd.String():
		return AppEnvProd
	default:
		return AppEnvProd
	}
}

func (cfg *AppConfig) IsEnvDebug() bool {
	return cfg.App.Env == AppEnvDebug
}

func (env AppEnv) String() string {
	return string(env)
}
