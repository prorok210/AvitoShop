package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockery --name=DB --output=./mocks --case=underscore
type DB interface {
	Acquire(ctx context.Context) (c *pgxpool.Conn, err error)
	AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
	AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
	Config() *pgxpool.Config
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Reset()
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Stat() *pgxpool.Stat
}

var DBConn DB

func ConnectDB() error {
	log.Println("Connecting to database")

	connUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))

	var err error
	ctx := context.Background()

	maxRetries := 3
	initialDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		DBConn, err = pgxpool.New(ctx, connUrl)
		if err == nil {
			if pingErr := DBConn.Ping(ctx); pingErr == nil {
				log.Println("Successfully connected to database")
				return nil
			}
			DBConn.Close()
		}

		log.Printf("Database not ready, retrying in %v... (%d/%d)", initialDelay, i+1, maxRetries)
		time.Sleep(initialDelay)
		initialDelay *= 2
	}

	log.Println("Failed to connect to database after retries")
	return errors.New("failed to connect to database")
}
