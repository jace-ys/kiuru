package crdb

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/slogger"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	DbName   string
	Insecure bool
}

type DB struct {
	*sqlx.DB
}

func Connect(c Config) (*DB, error) {
	connStr := fmt.Sprintf("postgresql://%s@%s:%d/%s?%s", c.User, c.Host, c.Port, c.DbName, sslMode(c.Insecure))
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		slogger.Info().Log("event", "crdb.connection_failure")
		return nil, errors.Wrap(err, "database connection failed")
	}
	slogger.Info().Log("event", "crdb.connection_success")
	return &DB{db}, nil
}

func sslMode(insecure bool) string {
	if insecure {
		return "sslmode=disable"
	}
	return "sslmode=require"
}

func (db *DB) Transact(fn func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
