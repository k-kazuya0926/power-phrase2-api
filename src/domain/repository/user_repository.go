package repository

import (
	"context"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

// UserRepository interface
type UserRepository interface {
	// TODO Contextを受け取るのはなぜ？
	Fetch(ctx context.Context) ([]*model.User, error)
	FetchByID(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
}
