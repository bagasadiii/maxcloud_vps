package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	dbUrl := os.Getenv("DATABASELINK")
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.New(context.Background(), dbUrl)
		if err == nil {
			if err = pool.Ping(context.Background()); err == nil {
				log.Println("Database connected successfully")
				break
			}
		}
		log.Printf("Retrying database connection... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil || pool == nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	table, err := os.ReadFile("table.sql")
	if err != nil {
		log.Fatalf("Failed to get query file: %v", err)
	}
	_, err = pool.Exec(context.Background(), string(table))
	if err != nil {
		log.Fatalf("Failed to create table query: %v", err)
	}
	return pool
}
