package gomodel

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	db *sql.DB
	DriverName string
	DataSourceName string
}

func (c Config) DB() (*sql.DB, error) {
	var err error = nil;
	
	if c.db == nil {
		c.db, err = sql.Open(c.DriverName, c.DataSourceName)
	}
	
	return c.db, err
}