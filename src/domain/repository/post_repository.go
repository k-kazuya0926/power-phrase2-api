package repository

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

type PostRepository interface {
	Create(post *model.Post) error
	Fetch(limit, page int, keyword string, userID int) (totalCount int, posts []*model.GetPostResult, err error)
	FetchByID(id int) (*model.GetPostResult, error)
	Update(post *model.Post) error
	Delete(id int) error
}
