package crdb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kru-travel/airdrop-go/pkg/slogger"

	_ "github.com/lib/pq"
)

type Config struct {
	Host          string
	Port          int
	User          string
	DbName        string
	RetryInterval time.Duration
	RetryCount    int
	Insecure      bool
}

type DB struct {
	*sqlx.DB
}

func NewCrdbClient(crdbConfig Config) (*DB, error) {
	db, err := connect(crdbConfig)
	if err != nil {
		slogger.Error().Log("event", "crdb_connection.failed", "msg", err)
		return nil, fmt.Errorf("database connection failed: %w", err)
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

func connect(c Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s@%s:%d/%s?%s", c.User, c.Host, c.Port, c.DbName, sslMode(c.Insecure))
	for i := 0; i < c.RetryCount; i++ {
		db, err := sqlx.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err != nil {
			slogger.Warn().Log("event", "crdb_connection.retry", "host", c.Host, "port", c.Port, "msg", err)
			time.Sleep(c.RetryInterval)
		} else {
			return db, nil
		}
	}
	return nil, fmt.Errorf("max number of retries exceeded")
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
