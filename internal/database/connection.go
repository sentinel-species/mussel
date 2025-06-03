package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"mussel/internal/types"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type connection struct {
	DB Pool

	readTimeout  time.Duration
	writeTimeout time.Duration
}

type Pool interface {
	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, arguments ...any) (rows pgx.Rows, err error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
}

func (c *connection) Ping(ctx context.Context) (err error) {
	pingCtx, pingCancel := context.WithTimeout(ctx, c.readTimeout)
	defer pingCancel()

	_, err = c.DB.Exec(pingCtx, "SELECT 1")
	if err != nil {
		fmt.Println("PING")
		return err
	}
	return nil
}

func (c *connection) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	return c.DB.Exec(ctx, sql, arguments...)
}

func (c *connection) Query(ctx context.Context, sql string, arguments ...any) (rows pgx.Rows, err error) {
	return c.DB.Query(ctx, sql, arguments...)
}

func (c *connection) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	return c.DB.QueryRow(ctx, sql, arguments...)
}

func Connect(cfg *types.Config) (pool Pool, err error) {
	ctx := context.Background()

	pgxConfig, err := pgxpool.ParseConfig(cfg.Database.Source)
	if err != nil {
		return nil, err
	}

	pgxConfig.MaxConnLifetime = cfg.Database.MaxConnectionLifetime
	pgxConfig.MaxConns = int32(cfg.Database.MaxConnections)
	pgxConfig.ConnConfig.ConnectTimeout = cfg.Database.ConnectTimeout

	connect := &connection{
		readTimeout:  cfg.Database.ReadTimeout,
		writeTimeout: cfg.Database.WriteTimeout,
	}

	connect.DB, err = pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for {
		err = connect.Ping(ctx)
		if err == nil || errors.Is(err, context.DeadlineExceeded) {
			break
		}
		time.Sleep(time.Second)
	}
	pool = connect.DB

	return
}
