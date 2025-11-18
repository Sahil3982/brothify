package database

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func RunMigrations(DB_URL string) {
	m, err := migrate.New("file://internal/migrations", DB_URL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
		return
	}
	log.Println("Database migrated successfully")

}
