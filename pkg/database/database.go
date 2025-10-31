package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectingDb() {
	godotenv.Load()

	connectionString := os.Getenv("DATABASE_URL")

	con, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer con.Close(context.Background())

	var greeting string
	err = con.QueryRow(context.Background(), "select 'Database Connection Successfully'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query Row Faild: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

}
