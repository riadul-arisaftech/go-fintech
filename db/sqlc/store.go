package db

import (
	"context"
	"database/sql"
	"fmt"
)

// # begin tx
// transfer money
// entry 1 in
// entry 2 out
// update balance
// # commit transaction

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fq func(q *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fq(q)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("encountered rollback error: %v", txErr)
		}
		return err
	}

	return tx.Commit()
}
