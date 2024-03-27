package databaze

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func DBConnect() {
	var err error
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"), os.Getenv("DB_JMENO"))
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Databaze se pokazila", err)
	}
}
