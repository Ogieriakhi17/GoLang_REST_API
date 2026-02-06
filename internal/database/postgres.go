package database

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string)(*pgxpool.Pool, error){
	var ctx context.Context = context.Background()

	// get config from the context pool
	var config *pgxpool.Config
	var err error
	config, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Println("Unable to Parse database URL: %v", err)
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil{
		log.Println("Unable to create connection pool")
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil {
		log.Println("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Yayy, successfully connected to Postgres database")
	return pool, nil
}