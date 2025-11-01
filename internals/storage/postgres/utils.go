package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func execInsertReturningID(ctx context.Context, db *sqlx.DB, query string, obj interface{}, idPtr *int) error {
	rows, err := db.NamedQueryContext(ctx, query, obj)
	if err != nil {
		return fmt.Errorf("exec insert with returning id failed: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(idPtr); err != nil {
			return fmt.Errorf("scan returned id failed: %w", err)
		}
		return nil
	}
	return fmt.Errorf("no id returned after insert")
}
