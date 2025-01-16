package services

import (
	"database/sql"
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
	return db, nil 
}


func checkDatabaseHealth(db *sql.DB) bool{
	return db.Ping() == nil
}
