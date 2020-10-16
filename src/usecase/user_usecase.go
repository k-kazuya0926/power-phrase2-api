package usecase

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// UserUseCase interfase
type UserUseCase interface {
	GetUsers() ([]*model.User, error)
	GetUser(id int) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int) error
}

type userUseCase struct {
	repository.UserRepository
}

// NewUserUseCase UserUseCaseを取得します.
func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}

func (usecase *userUseCase) GetUsers() ([]*model.User, error) {
	return usecase.UserRepository.Fetch()
}

func (usecase *userUseCase) GetUser(id int) (*model.User, error) {
	return usecase.UserRepository.FetchByID(id)
}

func (usecase *userUseCase) CreateUser(user *model.User) (*model.User, error) {
	return usecase.UserRepository.Create(user)
}

func (usecase *userUseCase) UpdateUser(user *model.User) (*model.User, error) {
	return usecase.UserRepository.Update(user)
}

func (usecase *userUseCase) DeleteUser(id int) error {
	return usecase.UserRepository.Delete(id)
}
