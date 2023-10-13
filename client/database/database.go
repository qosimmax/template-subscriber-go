// Package database contains a Postgres client and methods for communicating with the database.
package database

import (
	"context"
	"fmt"
	"template-subscriber-go/config"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

// Client holds the database client and prepared statements.
type Client struct {
	DB *sqlx.DB

	createAccountDataStmt     *sqlx.Stmt
	createCorpAccountDataStmt *sqlx.Stmt
	updateCorpAccountDataStmt *sqlx.Stmt
	deleteCorpAccountDataStmt *sqlx.Stmt
	getMytaxiIDStmt           *sqlx.Stmt
	payOrderStmt              *sqlx.Stmt
}

// Init sets up a new database client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseURL,
		config.DatabasePort,
		config.DatabaseDB,
		config.DatabaseOptions,
	)

	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(config.DatabaseMaxConnections)
	db.SetMaxIdleConns(config.DatabaseMaxIdleConnections)

	c.DB = db

	err = c.prepareStatements()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) prepareStatements() error {
	return nil
}

// Close closes the database connection and statements.
func (c *Client) Close() error {

	err := c.closeStatements()
	if err != nil {
		return err
	}

	err = c.DB.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	return nil
}

func (c *Client) closeStatements() error {

	return nil
}
