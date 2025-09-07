package database

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

func ConnectPostgreSQL(host, port, user, password, dbname string) *sql.DB {
    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    
    if err = db.Ping(); err != nil {
        log.Fatal("Failed to ping PostgreSQL:", err)
    }
    
    log.Println("✅ Connected to PostgreSQL")
    return db
}