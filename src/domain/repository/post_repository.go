package repository

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

type PostRepository interface {
	Create(post *model.Post) error
	FetchByID(id int) (*model.Post, error)
	Update(post *model.Post) error
	Delete(id int) error
}
