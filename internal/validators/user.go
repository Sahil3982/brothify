package validators

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UserExists(db *pgxpool.Pool, userID int) (bool,error){
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FORM users WHERE user_id = $1)`
	err := db.QueryRow(context.Background(),query,userID).Scan(&exists)
	return exists,err
}
