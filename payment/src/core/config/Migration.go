package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(pool *pgxpool.Pool) error {
	schema, err := os.ReadFile("init.sql")
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.Background(), string(schema))
	return err
}
