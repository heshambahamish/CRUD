package db

import (
	"database/sql"
	"log"
         "os"
	_ "github.com/lib/pq"
	
   
)
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = os.Getenv("DATABASE_URL")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database ping error:", err)
	}
}
