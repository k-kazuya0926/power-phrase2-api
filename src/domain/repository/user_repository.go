package repository

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

type UserRepository interface {
	Fetch() ([]*model.User, error)
	FetchByID(id int) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id int) error
}
