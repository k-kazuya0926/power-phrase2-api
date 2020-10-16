package datastore

import (
	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

type userRepository struct {
	connection *gorm.DB
}

// NewUserRepository UserRepositoryを取得します.
func NewUserRepository(connection *gorm.DB) repository.UserRepository {
	return &userRepository{connection}
}

func (repository *userRepository) Fetch() ([]*model.User, error) {
	var (
		users []*model.User
		err   error
	)
	err = repository.connection.Order("id desc").Find(&users).Error
	return users, err
}

func (repository *userRepository) FetchByID(id int) (*model.User, error) {
	u := &model.User{ID: id}
	err := repository.connection.First(u).Error
	return u, err
}

func (repository *userRepository) Create(u *model.User) (*model.User, error) {
	err := repository.connection.Create(u).Error
	return u, err
}

func (repository *userRepository) Update(u *model.User) (*model.User, error) {
	err := repository.connection.Model(u).Update(u).Error
	return u, err
}

func (repository *userRepository) Delete(id int) error {
	u := &model.User{ID: id}
	err := repository.connection.Delete(u).Error
	return err
}
