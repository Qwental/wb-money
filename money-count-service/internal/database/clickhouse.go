package database

import (
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
)

func NewClickHouseDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("clickhouse", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
