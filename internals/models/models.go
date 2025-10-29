package models

import (
	"time"
)

type Post struct {
	ID            int       `db:"id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	Author        *string   `db:"author"`
	AllowComments bool      `db:"allow_comments"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type PostWithComments struct {
	Post
	Comments []Comment
}

type Comment struct {
	ID        int       `db:"id"`
	PostID    int       `db:"post_id"`
	ParentID  *int      `db:"parent_id"`
	Author    *string   `db:"author"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
