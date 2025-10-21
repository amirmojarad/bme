package database

import (
	"bme/conf"
	"bme/pkg/errorext"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToPostgres(cfg *conf.AppConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.SslMode,
		cfg.Postgres.Timezone,
	)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errorext.New(err, errorext.ErrGeneralOccurrence)
	}

	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(cfg.Postgres.ConnectionMaxLifetime)

	return sqlDB, nil
}
