package databaze

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func DBConnect() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error .env", err)
	}
	connStr := fmt.Sprintf("postgresql://%s:%s@%s", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"))
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Databaze se pokazila", err)
	}
}
