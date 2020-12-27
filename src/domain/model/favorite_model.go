// Package model Domain Model
package model

import (
	"time"
)

// Favorite favoritesテーブルに対応する構造体。
type Favorite struct {
	ID        int       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:current_timestamp"`
	PostID    int       `json:"post_id" gorm:"not null;default:0"`
	UserID    int       `json:"user_id" gorm:"not null;default:0"`
}
