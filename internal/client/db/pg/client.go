package pg

import (
	"context"

	"github.com/Ippolid/chat-server/internal/client/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC db.DB
}

// New создает новый экземпляр pgClient
func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to db: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
