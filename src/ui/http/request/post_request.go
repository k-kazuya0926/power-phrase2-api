// Package request リクエストを表す構造体を定義
package request

type (
	// CreatePostRequest 投稿登録リクエスト
	CreatePostRequest struct {
		UserID   int    `json:"user_id" validate:"required,min=1"`
		Title    string `validate:"required,max=100"`
		Speaker  string `validate:"required,max=100"`
		Detail   string `validate:"max=500"`
		MovieURL string `json:"movie_url" validate:"max=200"`
	}

	// GetPostsRequest 投稿一覧取得リクエスト
	GetPostsRequest struct {
		Limit       int    `json:"limit" validate:"required,min=1"`
		Page        int    `json:"page" validate:"required,min=1"`
		Keyword     string `json:"keyword" validate:"max=100"`
		PostUserID  int    `json:"post_user_id" validate:"min=0"`
		LoginUserID int    `json:"login_user_id" validate:"min=0"`
	}

	// GetPostRequest 投稿詳細取得リクエスト
	GetPostRequest struct {
		ID          int `validate:"min=1"`
		LoginUserID int `json:"login_user_id" validate:"min=0"`
	}

	// UpdatePostRequest 投稿更新リクエスト
	UpdatePostRequest struct {
		ID       int    `validate:"required,min=1"`
		Title    string `validate:"required,max=100"`
		Speaker  string `validate:"required,max=100"`
		Detail   string `validate:"max=500"`
		MovieURL string `json:"movie_url" validate:"max=200"`
	}

	// DeletePostRequest 投稿削除リクエスト
	DeletePostRequest struct {
		ID int `validate:"min=1"`
	}

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
		UserID int `json:"user_id" validate:"required,min=1"`
		PostID int `json:"post_id" validate:"required,min=1"`
	}
)
