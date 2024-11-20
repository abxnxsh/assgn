package database

import (
    "database/sql"
    "fmt"

    _ "github.com/lib/pq" 
)

var DB *sql.DB

func ConnectDB() error {
    connStr := "host=localhost port=5432 user=product_user password=password dbname=product_db sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }

    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    DB = db
    fmt.Println("Database connection successful!")
    return nil
}
