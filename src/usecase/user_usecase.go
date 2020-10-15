package usecase

import (
	"context"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// UserUseCase interfase
type UserUseCase interface {
	// TODO Modelを返していいのか？
	// TODO Contextを受け取るのはなぜ？
	GetUsers(ctx context.Context) ([]*model.User, error)
	GetUser(ctx context.Context, id int) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userUseCase struct {
	repository.UserRepository
}

// NewUserUseCase UserUseCaseを取得します.
func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}

func (usecase *userUseCase) GetUsers(ctx context.Context) ([]*model.User, error) {
	return usecase.UserRepository.Fetch(ctx)
}

func (usecase *userUseCase) GetUser(ctx context.Context, id int) (*model.User, error) {
	return usecase.UserRepository.FetchByID(ctx, id)
}

func (usecase *userUseCase) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return usecase.UserRepository.Create(ctx, user)
}

func (usecase *userUseCase) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	return usecase.UserRepository.Update(ctx, user)
}

func (usecase *userUseCase) DeleteUser(ctx context.Context, id int) error {
	return usecase.UserRepository.Delete(ctx, id)
}
