package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBConn *pgxpool.Pool

func ConnectDB(ctx context.Context, connURL string) error {
	log.Println("Connecting to database")

	conn, err := pgxpool.New(ctx, connURL)
	if err != nil || conn == nil {
		log.Println("Connection failed")
		return err
	}

	DBConn = conn
	if err = DBConn.Ping(ctx); err != nil {
		log.Println("Connection failed")
		return err
	}

	log.Println("Successfully connected to database")
	return nil
}
