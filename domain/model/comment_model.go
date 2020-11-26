// Package model Domain Model
package model

import (
	"time"
)

// Comment commentsテーブルに対応する構造体。
type Comment struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at"`
	PostID    int        `json:"post_id" gorm:"not null;default:0"`
	UserID    int        `json:"user_id" gorm:"not null;default:0"`
	Body      string     `json:"body" gorm:"type:varchar(256);not null;default:''"`
}
