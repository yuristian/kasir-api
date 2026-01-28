package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25) //define max open connection di dalam 1 waktu.
	db.SetMaxIdleConns(5)  //kalo ga ada transaksi, 5 connection yang tersedia / yg ready.

	log.Println("Database connected successfully")
	return db, nil
}
