package services

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresDb(address string) (*sql.DB, error) {
	db, err := sql.Open("postgres", address)
	if err != nil{
		return nil, err
	}
	if !checkDatabaseHealth(db){
		return nil, err 
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, nil 
}


func checkDatabaseHealth(db *sql.DB) bool{
	return db.Ping() == nil
}
