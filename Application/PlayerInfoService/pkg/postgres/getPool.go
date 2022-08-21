package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const maxAttempts = 10

func GetConnectionPool(connectionString string) *pgxpool.Pool {
	ctx := context.Background()
	connectionPool, _ := NewClient(ctx, maxAttempts, connectionString)

	return connectionPool
}

type PoolInterface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

func NewClient(ctx context.Context, maxAttempts int, connectionString string) (connectionPool *pgxpool.Pool, err error) {
	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		connectionPool, err = pgxpool.Connect(ctx, connectionString)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return connectionPool, nil
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return err
}
