package database

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type TxCtxKey struct{}

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error
}

type GormWrapper struct {
	db *gorm.DB
}

func NewGormWrapper(db *gorm.DB) *GormWrapper {
	return &GormWrapper{
		db: db,
	}
}

func (g GormWrapper) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(TxCtxKey{}).(*gorm.DB)
	if ok {
		return tx
	}

	return g.db.WithContext(ctx)
}

func (g GormWrapper) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...*sql.TxOptions) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		newCtx := context.WithValue(ctx, TxCtxKey{}, tx)

		return fn(newCtx)
	}, opts...)
}
