package database

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5"
)

func main() {
	connectionString := "postgres://username:Sahil@2002@localhost:5432/brothifydb"
	con, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)		
	}

	defer con.Close(context.Background())
	

}
