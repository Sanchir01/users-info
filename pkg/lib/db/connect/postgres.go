package connect

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Sanchir01/candles_backend/pkg/lib/utils"
	"github.com/Sanchir01/users-info/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PGXNew(cfg *config.Config, ctx context.Context) (*pgxpool.Pool, error) {
	var dsn string

	switch cfg.Env {
	case "development":
		dsn = fmt.Sprintf(
			"postgresql://postgres:postgres@localhost:5439/test",
		)
	case "production":
		dsn = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s",
			cfg.DB.User, os.Getenv("DB_PASSWORD_PROD"),
			cfg.DB.Host, cfg.DB.Port, cfg.DB.Database,
		)
	}
	var pool *pgxpool.Pool
	var err error

	err = utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, cfg.DB.MaxAttempts, 5*time.Second)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
