package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectingDb() *pgxpool.Pool {
	godotenv.Load()
	connectionString := os.Getenv("DATABASE_URL")
	// RunMigrations(connectionString)

	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = pool.QueryRow(context.Background(), "select 'Database Connection Successfully...'").Scan(&greeting)
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
	DB = pool
	return pool
}
