// Package datastore Infra層のリポジトリ
package datastore

import (
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// commentRepository 構造体
type commentRepository struct {
}

// NewCommentRepository CommentRepositoryを生成する。
func NewCommentRepository() repository.CommentRepository {
	return &commentRepository{}
}

// Create 登録
func (repository *commentRepository) Create(comment *model.Comment) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(comment).Error
}

// Fetch 一覧取得
func (repository *commentRepository) Fetch(postID, limit, page int) (totalCount int, comments []*model.GetCommentResult, err error) {
	db := conf.NewDBConnection()
	defer db.Close()

	countDb := conf.NewDBConnection()
	defer countDb.Close()

	offset := limit * (page - 1)

	if err = countDb.Model(&model.Comment{}).Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	if err = db.Table("comments").
		Select("comments.*, users.name as user_name, users.image_file_path as user_image_file_path").
		Joins("LEFT JOIN users on users.id = comments.user_id AND users.deleted_at IS NULL").
		Order("id DESC").Limit(limit).Offset(offset).
		Find(&comments).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, comments, err
}

// Delete 削除
func (repository *commentRepository) Delete(id int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	comment := model.Comment{ID: id}
	return db.Delete(&comment).Error
}
