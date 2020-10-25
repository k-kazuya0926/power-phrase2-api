package datastore

import (
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

type postRepository struct {
}

// NewPostRepository PostRepositoryを取得します.
func NewPostRepository() repository.PostRepository {
	return &postRepository{}
}

func (repository *postRepository) Create(post *model.Post) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	return connection.Create(post).Error
}

func (repository *postRepository) FetchByID(id int) (*model.Post, error) {
	connection := conf.NewDBConnection()
	defer connection.Close()

	u := model.Post{ID: id}
	err := connection.First(&u).Error

	return &u, err
}

func (repository *postRepository) Update(u *model.Post) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	return connection.Model(u).Update(u).Error
}

func (repository *postRepository) Delete(id int) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	post := model.Post{ID: id}
	return connection.Delete(&post).Error
}
