package models

import "time"

type Post struct {
	ID            int       `db:"id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	UserID        int       `db:"user_id"`
	AllowComments bool      `db:"allow_comments"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
