package request

type (
	CreatePostRequest struct {
		UserID   int    `json:"user_id" validate:"required,min=1"`
		Title    string `validate:"required,max=100"`
		Speaker  string `validate:"required,max=100"`
		Detail   string `validate:"max=500"`
		MovieURL string `json:"movie_url" validate:"max=200"`
	}

	GetPostsRequest struct {
		Limit   int    `json:"limit" validate:"required,min=1"`
		Page    int    `json:"page" validate:"required,min=1"`
		Keyword string `json:"keyword" validate:"max=100"`
		UserID  int    `json:"user_id" validate:"min=0"`
	}

	GetPostRequest struct {
		ID int `validate:"min=1"`
	}

	UpdatePostRequest struct {
		ID       int    `validate:"required,min=1"`
		Title    string `validate:"required,max=100"`
		Speaker  string `validate:"required,max=100"`
		Detail   string `validate:"max=500"`
		MovieURL string `json:"movie_url" validate:"max=200"`
	}

	DeletePostRequest struct {
		ID int `validate:"min=1"`
	}
)
