package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigration(db *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
   			 user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   			 name VARCHAR(255) NOT NULL,
  			 email VARCHAR(255) UNIQUE NOT NULL,
   			 password TEXT NOT NULL,
   			 created_at TIMESTAMP DEFAULT NOW(),
   			 updated_at TIMESTAMP DEFAULT NOW()
		);
		`,
		`CREATE TABLE IF NOT EXISTS dishes (
			dish_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    		dish_name VARCHAR(255) NOT NULL,
    		cat_id INT,
    		price NUMERIC(10,2) NOT NULL,
    		description TEXT,
    		dish_url TEXT, 
    		availability BOOLEAN DEFAULT TRUE,
    		rating NUMERIC(3,2) DEFAULT 0,
    		highlight BOOLEAN DEFAULT FALSE,
    		created_at TIMESTAMP DEFAULT NOW(),
    		updated_at TIMESTAMP DEFAULT NOW()
		);
`,
		`CREATE TABLE IF NOT EXISTS reservations (
			reservation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id INT NOT NULL REFERENCES users(user_id),
			table_number INT NOT NULL,
			reservation_person_name VARCHAR(255),
			reservation_person_email VARCHAR(255),
			reservation_person_mobile_number VARCHAR(20),
			reservation_time VARCHAR(50),
			reservation_date VARCHAR(50),
			number_of_guests INT,
			special_requests TEXT,
			status VARCHAR(20), 
			amount NUMERIC(10,2),
			payment_id VARCHAR(255),
			signature VARCHAR(255),
			payment_status VARCHAR(50),
			invoice_url TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);`,
		` CREATE TABLE IF NOT EXISTS reservation_dishes (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			reservation_id INT NOT NULL REFERENCES reservations(reservation_id) ON DELETE CASCADE,
			dish_id UUID NOT NULL REFERENCES dishes(dish_id),
			UNIQUE(reservation_id, dish_id)
		);
		`,
	}

	for _, q := range queries {
		_, err := db.Exec(context.Background(), q)
		if err != nil {
			fmt.Println("Migration failed:", err)
			return err
		}		
	}

	fmt.Println("✅ Migrations applied successfully")
	return nil

}
