// Package datastore Infra層のリポジトリ
package datastore

import (
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// postRepository 構造体
type postRepository struct {
}

// NewPostRepository PostRepositoryを生成する。
func NewPostRepository() repository.PostRepository {
	return &postRepository{}
}

// Create 登録
func (repository *postRepository) Create(post *model.Post) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(post).Error
}

// Fetch 一覧取得。キーワード検索を行わない場合はkeywordに空文字を指定する。ユーザーを限定しない場合はuserIDに0を指定する。
func (repository *postRepository) Fetch(limit, page int, keyword string, userID int) (totalCount int, posts []*model.GetPostResult, err error) {
	db := conf.NewDBConnection()
	defer db.Close()

	countDb := conf.NewDBConnection()
	defer countDb.Close()

	offset := limit * (page - 1)

	if keyword != "" {
		countDb = countDb.Where("title LIKE ?", "%"+keyword+"%").Or("speaker LIKE ?", "%"+keyword+"%").Or("detail LIKE ?", "%"+keyword+"%")
		db = db.Where("title LIKE ?", "%"+keyword+"%").Or("speaker LIKE ?", "%"+keyword+"%").Or("detail LIKE ?", "%"+keyword+"%")
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
		Joins("JOIN users on users.id = posts.user_id AND users.deleted_at IS NULL").
		Order("id DESC").Limit(limit).Offset(offset).
		Find(&posts).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, posts, err
}

// FetchByID 1件取得
func (repository *postRepository) FetchByID(id int) (*model.GetPostResult, error) {
	db := conf.NewDBConnection()
	defer db.Close()

	post := model.GetPostResult{}
	post.ID = id
	if err := db.Table("posts").
		Select("posts.*, users.name as user_name, users.image_file_path as user_image_file_path").
		Joins("JOIN users on users.id = posts.user_id AND users.deleted_at IS NULL").
		First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// Update 更新
func (repository *postRepository) Update(u *model.Post) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Model(u).Update(u).Error
}

// Delete 削除
func (repository *postRepository) Delete(id int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	post := model.Post{ID: id}
	return db.Delete(&post).Error
}
