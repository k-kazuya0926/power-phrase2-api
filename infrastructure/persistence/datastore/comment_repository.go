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
