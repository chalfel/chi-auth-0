package db

import (
	"context"
)

type Db struct{}

func (d *Db) Exec(ctx context.Context, sql string, args ...string) error {
	return nil
}

func NewDb() *Db {
	return &Db{}
}
