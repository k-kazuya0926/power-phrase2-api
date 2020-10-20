package request

type (
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
)
