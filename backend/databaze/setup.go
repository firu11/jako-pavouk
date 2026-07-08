package databaze

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(connStr string) {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to database")
}

func Close() {
	_ = DB.Close()
}
