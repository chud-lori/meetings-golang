package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

    _ "github.com/lib/pq"
)

func NewPostgreDB() *sql.DB {

    const (
        host = "localhost"
        port = 5432
        user = "postgres"
        password = "root"
        dbname = "jub_dup"
    )

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)

    if err != nil {
        panic(err)
    }

    db.SetMaxIdleConns(5)
    db.SetMaxOpenConns(20)
    db.SetConnMaxLifetime(60 * time.Minute)
    db.SetConnMaxIdleTime(10 * time.Minute)

    return db
}

