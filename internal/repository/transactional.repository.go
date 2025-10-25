package repository

import (
	"bme/database"
	"context"
	"github.com/pkg/errors"
)

type Transactional struct {
	*database.GormWrapper
}

func NewTransactional(db *database.GormWrapper) *Transactional {
	return &Transactional{
		db,
	}
}

func (r Transactional) BeginTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return errors.WithStack(r.WithTransaction(ctx, fn))
}
