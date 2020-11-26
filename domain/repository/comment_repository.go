// Package repository Domain Service層のリポジトリ
package repository

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

// CommentRepository commentsテーブルへのアクセスを行うインターフェース。
type CommentRepository interface {
	Create(comment *model.Comment) error
}
