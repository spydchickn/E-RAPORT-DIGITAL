package config

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
    var err error
    // Enable parseTime for DATE/TIMESTAMP scanning, set UTF8MB4 charset and local time zone
    dsn := "root:@tcp(127.0.0.1:3306)/e_raport?parseTime=true&charset=utf8mb4&loc=Local"
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("sql.Open:", err)
    }
    if err := DB.Ping(); err != nil {
        log.Fatal("db.Ping:", err)
    }
    fmt.Println("Database connected!")
}
