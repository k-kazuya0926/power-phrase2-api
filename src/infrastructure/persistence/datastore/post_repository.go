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

func (repository *postRepository) Fetch(limit, page int, keyword string, userID int) (totalCount int, posts []*model.GetPostResult, err error) {
	db := conf.NewDBConnection()
	defer db.Close()

	countDb := conf.NewDBConnection()
	defer countDb.Close()

	offset := limit * (page - 1)

	if keyword != "" {
		// TODO タイトル以外も対象にする
		countDb = countDb.Where("title LIKE ?", "%"+keyword+"%")
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}

	if userID > 0 { // ユーザーIDが指定されている場合
		countDb = countDb.Where("user_id = ?", userID)
		db = db.Where("user_id = ?", userID)
	}

	if err = countDb.Model(&model.Post{}).Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	if err = db.Table("posts").
		Select("posts.*, users.name as user_name, users.image_file_path as user_image_file_path").
		Joins("LEFT JOIN users on users.id = posts.user_id AND users.deleted_at IS NULL").
		Order("id DESC").Limit(limit).Offset(offset).
		Scan(&posts).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, posts, err
}

func (repository *postRepository) FetchByID(id int) (*model.GetPostResult, error) {
	db := conf.NewDBConnection()
	defer db.Close()

	post := model.GetPostResult{}
	post.ID = id
	if err := db.Table("posts").
		Select("posts.*, users.name as user_name, users.image_file_path as user_image_file_path").
		Joins("LEFT JOIN users on users.id = posts.user_id AND users.deleted_at IS NULL").
		First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
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
