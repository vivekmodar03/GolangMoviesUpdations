package db

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", "root:system@0987@tcp(127.0.0.1:3305)/go_movie_db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Ping failed:", err)
	}

	log.Println("MySQL Connected!")
}
