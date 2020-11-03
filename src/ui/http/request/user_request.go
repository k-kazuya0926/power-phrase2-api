package request

type (
	CreateUserRequest struct {
		Name     string `json:"name" validate:"required,max=50"`
		Email    string `json:"email" validate:"required,email,max=100"`
		Password string `json:"password" validate:"required,max=100"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email,max=100"`
		Password string `json:"password" validate:"required,max=100"`
	}

	GetUserRequest struct {
		ID int `json:"id" validate:"min=1"`
	}

	UpdateUserRequest struct {
		ID       int    `json:"id" validate:"required,min=1"`
		Name     string `json:"name" validate:"required,max=50"`
		Email    string `json:"email" validate:"required,email,max=100"`
		Password string `json:"password" validate:"max=100"`
		ImageURL string `json:"image_url" validate:"max=100"` // TODO 削除
	}

	DeleteUserRequest struct {
		ID int `json:"id" validate:"min=1"`
	}
)
