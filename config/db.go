package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dsn := "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable"
	DB, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	fmt.Println("✅ Database connected")
}
