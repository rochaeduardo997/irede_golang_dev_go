package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

func NewDatabaseConnection() (result *sql.DB, err error) {
	DRIVER := "mysql"
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_DB")

	cfg := mysql.Config{
		User:                 USER,
		Passwd:               PASSWORD,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", HOST, PORT),
		DBName:               DBNAME,
		AllowNativePasswords: true,
	}
	result, err = sql.Open(DRIVER, cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return
}
