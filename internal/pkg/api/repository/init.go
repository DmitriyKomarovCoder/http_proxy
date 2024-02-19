package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func NewPostgresConnect(ctx context.Context, host, user, password, name string, port int) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)

	db, err := pgxpool.New(ctx, connStr)

	err = db.Ping(context.Background())
	if err != nil {
		log.Println("Error while ping to DB", err)
		return nil, err
	}

	return db, nil
}
