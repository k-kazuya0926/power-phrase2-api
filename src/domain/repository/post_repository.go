// Package repository Domain Service層のリポジトリ
package repository

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

// PostRepository postsや関連テーブルへのアクセスを行うインターフェース。
type PostRepository interface {
	// 投稿登録
	Create(post *model.Post) error
	// 投稿一覧取得
	Fetch(limit, page int, keyword string, userID int) (totalCount int, posts []*model.GetPostResult, err error)
	// 投稿詳細取得
	FetchByID(id int) (*model.GetPostResult, error)
	// 投稿更新
	Update(post *model.Post) error
	// 投稿削除
	Delete(id int) error

	// コメント登録
	CreateComment(comment *model.Comment) error
	// コメント一覧取得
	FetchComments(postID, limit, page int) (totalCount int, comments []*model.GetCommentResult, err error)
	// 投稿削除
	DeleteComment(id int) error

	// お気に入り登録
	CreateFavorite(favorite *model.Favorite) error
	// お気に入り一覧取得
	FetchFavorites(userID, limit, page int) (totalCount int, posts []*model.GetPostResult, err error)
	// お気に入り削除
	DeleteFavorite(id int) error
}
