package database

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func ConnectingDb() *pgx.Conn {
	godotenv.Load()

	connectionString := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Database Connection Successfully...'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Row Faild: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
	DB = conn
	return conn

}
