// Package request リクエストを表す構造体を定義
package request

type (
	// CreateCommentRequest コメント登録リクエスト
	CreateCommentRequest struct {
		PostID int    `json:"post_id" validate:"required,min=1"`
		UserID int    `json:"user_id" validate:"required,min=1"`
		Body   string `validate:"required,max=200"`
	}

	// GetCommentsRequest コメント一覧取得リクエスト
	GetCommentsRequest struct {
		PostID int `json:"post_id" validate:"required,min=1"`
		Limit  int `json:"limit" validate:"required,min=1"`
		Page   int `json:"page" validate:"required,min=1"`
	}

	// DeleteCommentRequest コメント削除リクエスト
	DeleteCommentRequest struct {
		ID int `validate:"min=1"`
	}
)
