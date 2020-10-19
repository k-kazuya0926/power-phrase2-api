package model

import (
	"time"
)

// User struct
type User struct {
	ID        int        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null;default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" gorm:"type:varchar(256);not null;default:''" validate:"required,max=50"`
	Email     string     `json:"email" gorm:"type:varchar(256);not null;default:'';unique" validate:"required,email,max=100"`
	Password  string     `json:"password" gorm:"type:varchar(256);not null;default:''" validate:"required,max=100"`
	ImageURL  string     `json:"image_url" gorm:"type:varchar(256);not null;default:''" validate:"max=100"`
}
