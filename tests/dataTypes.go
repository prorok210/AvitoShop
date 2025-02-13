package tests

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type FakeUserRow struct {
	UserID   int
	Username string
	Password string
	Err      error
}

func (r *FakeUserRow) Scan(dest ...interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	if len(dest) >= 1 {
		if ref, ok := dest[0].(*int); ok {
			*ref = r.UserID
		}
	}
	if len(dest) >= 2 {
		if ref, ok := dest[1].(*string); ok {
			*ref = r.Password
		}
	}

	if len(dest) >= 3 {
		if ref, ok := dest[2].(*string); ok {
			*ref = r.Username
		}
	}
	return nil
}

type FakeMerchRowInt struct {
	MerchID int
	Price   int
	Err     error
}

func (r *FakeMerchRowInt) Scan(dest ...interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	if len(dest) >= 1 {
		if ref, ok := dest[0].(*int); ok {
			*ref = r.MerchID
		}
	}
	if len(dest) >= 2 {
		if ref, ok := dest[1].(*int); ok {
			*ref = r.Price
		}
	}
	return nil
}

type Tx struct {
	mock.Mock
	pgx.Tx
}

func (tx *Tx) Begin(ctx context.Context) (pgx.Tx, error) {
	return tx, nil
}

func (tx *Tx) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	allArgs := []interface{}{ctx, query}
	allArgs = append(allArgs, args...)
	arguments := tx.Called(allArgs...)
	return arguments.Get(0).(pgconn.CommandTag), arguments.Error(1)
}

func (tx *Tx) Commit(ctx context.Context) error {
	arguments := tx.Called(ctx)
	return arguments.Error(0)
}

func (tx *Tx) Rollback(ctx context.Context) error {
	arguments := tx.Called(ctx)
	return arguments.Error(0)
}
