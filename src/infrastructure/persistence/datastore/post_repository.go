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
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(post).Error
}

func (repository *postRepository) Fetch(limit, page int, keyword string) ([]*model.Post, error) {
	db := conf.NewDBConnection()
	defer db.Close()

	offset := limit * (page - 1)

	if keyword != "" {
		// TODO タイトル以外も対象にする
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	var posts []*model.Post
	err := db.Order("id desc").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, err
}

func (repository *postRepository) FetchByID(id int) (*model.Post, error) {
	db := conf.NewDBConnection()
	defer db.Close()

	u := model.Post{ID: id}
	err := db.First(&u).Error

	return &u, err
}

func (repository *postRepository) Update(u *model.Post) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Model(u).Update(u).Error
}

func (repository *postRepository) Delete(id int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	post := model.Post{ID: id}
	return db.Delete(&post).Error
}
