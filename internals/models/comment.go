package models

import "time"

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	ParentID  *int      `db:"parent_id"`
	UserID    int       `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
