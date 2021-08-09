package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/taudelta/notifyhub/internal/notifyhub/config"
)

var db *sql.DB

func CreateDefaultConnection(connString string) error {
	var err error
	db, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetConnString(dbConfig config.DatabaseConfig) string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.SSLMode,
	)
}

func Connection() *sql.DB {
	return db
}
