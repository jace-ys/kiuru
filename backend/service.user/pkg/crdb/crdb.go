package crdb

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

func NewCrdbClient(c Config) (*DB, error) {
	connStr := fmt.Sprintf("postgresql://%s@%s:%d/%s?%s", c.User, c.Host, c.Port, c.DbName, sslMode(c.Insecure))
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		slogger.Info().Log("event", "crdb_connection.failed", "msg", err)
		return nil, errors.Wrap(err, "database connection failed")
	}
	db.MapperFunc(toLowerSnakeCase)
	slogger.Info().Log("event", "crdb_connection.success")
	return &DB{db}, nil
}

func (db *DB) Transact(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
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

func sslMode(insecure bool) string {
	if insecure {
		return "sslmode=disable"
	}
	return "sslmode=require"
}

func toLowerSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
