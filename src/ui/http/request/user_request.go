package request

type (
	// TODO DB存在チェックはできるか？

	CreateUserRequest struct {
		Name     string `validate:"required,max=50"`
		Email    string `validate:"required,email,max=100"`
		Password string `validate:"required,max=100"`
		ImageURL string `validate:"max=100"`
	}

	LoginRequest struct {
		Email    string `validate:"required,email,max=100"`
		Password string `validate:"required,max=100"`
	}

	GetUserRequest struct {
		UserID int `validate:"min=1"`
	}

	UpdateUserRequest struct {
		UserID   int    `validate:"required,min=1"`
		Name     string `validate:"max=50"`
		Email    string `validate:"max=100"` // TODO 空でない場合、形式チェック
		Password string `validate:"max=100"`
		ImageURL string `validate:"max=100"`
	}

	DeleteUserRequest struct {
		UserID int `validate:"min=1"`
	}
)
