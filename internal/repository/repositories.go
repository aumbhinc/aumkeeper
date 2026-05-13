package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repositories {
	return &Repositories{DB: db}
}