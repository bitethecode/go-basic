package db

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	host     = "192.168.1.157"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "db"
)

func OpenDb() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(1 * time.Second)
	db.SetConnMaxLifetime(30 * time.Second)

	return db
}
