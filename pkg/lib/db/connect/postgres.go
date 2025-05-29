package connect

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Sanchir01/users-info/pkg/lib/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PGXNew(ctx context.Context, user, host, db, port string, maxAttempts int) (*pgxpool.Pool, error) {
	//todo:: delete ssl disable for production
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		user, os.Getenv("DB_PASSWORD"),
		host, port, db,
	)

	var pool *pgxpool.Pool

	err := utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var err error
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
