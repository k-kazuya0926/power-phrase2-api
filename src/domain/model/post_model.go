package model

import (
	"time"
)

// Post struct
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
