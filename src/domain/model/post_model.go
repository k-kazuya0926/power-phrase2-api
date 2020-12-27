// Package model Domain Model
package model

import (
	"time"
)

// Post postsテーブルに対応する構造体。
type Post struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    int        `json:"user_id" gorm:"not null;default:0"`
	Title     string     `json:"title" gorm:"type:varchar(256);not null;default:''"`
	Speaker   string     `json:"speaker" gorm:"type:varchar(256);not null;default:''"`
	Detail    string     `json:"detail" gorm:"type:varchar(512);not null;default:''"`
	MovieURL  string     `json:"movie_url" gorm:"type:varchar(256);not null;default:''"`
}

// GetPostResult GetPostの戻り値として使用される構造体。
type GetPostResult struct {
	Post
	EmbedMovieURL     string `json:"embed_movie_url"`
	UserName          string `json:"user_name"`
	UserImageFilePath string `json:"user_image_file_path"`
	CommentCount      int    `json:"comment_count"`
	IsFavorite        bool   `json:"is_favorite"`
}

// Favorite favoritesテーブルに対応する構造体。
type Favorite struct {
	ID        int       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:current_timestamp"`
	UserID    int       `json:"user_id" gorm:"not null;default:0"`
	PostID    int       `json:"post_id" gorm:"not null;default:0"`
}
