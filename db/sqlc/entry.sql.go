// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: entry.sql

package db

import (
	"context"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entries, error) {
	row := q.db.QueryRowContext(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entries
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 AND account_id = $2 
LIMIT 1
`

type GetEntryParams struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
}

func (q *Queries) GetEntry(ctx context.Context, arg GetEntryParams) (Entries, error) {
	row := q.db.QueryRowContext(ctx, getEntry, arg.ID, arg.AccountID)
	var i Entries
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $3
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListEntriesParams struct {
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
	AccountID int64 `json:"account_id"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entries, error) {
	rows, err := q.db.QueryContext(ctx, listEntries, arg.Limit, arg.Offset, arg.AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entries{}
	for rows.Next() {
		var i Entries
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
