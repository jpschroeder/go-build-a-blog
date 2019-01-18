package main

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/acme/autocert"
)

type SqliteCache struct {
	db *sql.DB
}

func (c SqliteCache) Get(ctx context.Context, name string) ([]byte, error) {
	sql := `
		select Data from cache where Name = ?
	`
	row := c.db.QueryRow(sql, name)

	var data []byte
	err := row.Scan(&data)
	if err != nil {
		return nil, autocert.ErrCacheMiss
	}
	return data, nil
}

func (c SqliteCache) Put(ctx context.Context, name string, data []byte) error {
	sql := `
		insert into cache(Name, Data) values(?, ?)
		on conflict(Name) do update set Data=excluded.Data
	`
	_, err := c.db.Exec(sql, name, data)
	return err
}

func (c SqliteCache) Delete(ctx context.Context, name string) error {
	sql := `
		delete from cache where Name = ?
	`
	_, err := c.db.Exec(sql, name)
	return err
}
