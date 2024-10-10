package db

import (
	"database/sql"
	"fmt"

  _ "github.com/lib/pq"
)


func NewPGStorage(host, user, dbname, password, sslmode string, port int64) (*sql.DB, error) {
  connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s port=%d", host, user, dbname, password, sslmode, port)

  db, err := sql.Open("postgres", connStr)
  if err != nil {
    return nil, err
  }

  if err := db.Ping(); err != nil {
    return nil, err
  }

  return db, nil
} 
