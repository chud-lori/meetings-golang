package infrastructure

import (
	"database/sql"
	"fmt"
	"time"
    "os"

    _ "github.com/lib/pq"
)

func NewPostgreDB() *sql.DB {

    const (
    //    host = os.Getenv("PSQL_HOST")
        port = 5432
    //    user = ""
    //    password = ""
        dbname = "jub_dup"
    )

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PSQL_HOST"), port, os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASSWORD"), dbname)
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

