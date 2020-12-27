// Package request リクエストを表す構造体を定義
package request

type (
	// CreateFavoriteRequest お気に入り登録リクエスト
	CreateFavoriteRequest struct {
		UserID int `json:"user_id" validate:"required,min=1"`
		PostID int `json:"post_id" validate:"required,min=1"`
	}

	// GetFavoritesRequest お気に入り一覧取得リクエスト
	GetFavoritesRequest struct {
		UserID int `json:"user_id" validate:"required,min=1"`
		Limit  int `json:"limit" validate:"required,min=1"`
		Page   int `json:"page" validate:"required,min=1"`
	}

	// DeleteFavoriteRequest お気に入り削除リクエスト
	DeleteFavoriteRequest struct {
		ID int `json:"id" validate:"min=1"`
	}
)
