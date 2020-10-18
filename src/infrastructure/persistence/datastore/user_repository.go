package datastore

import (
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

type userRepository struct {
}

// NewUserRepository UserRepositoryを取得します.
func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (repository *userRepository) Create(u *model.User) (*model.User, error) {
	connection := conf.NewDBConnection()
	defer connection.Close()

	err := connection.Create(u).Error
	u.Password = ""
	return u, err
}

func (repository *userRepository) FetchByEmail(email string) (*model.User, error) {
	var u model.User

	connection := conf.NewDBConnection()
	defer connection.Close()

	err := connection.Where("email = ?", email).First(&u).Error
	return &u, err
}

func (repository *userRepository) Fetch() ([]*model.User, error) {
	connection := conf.NewDBConnection()
	defer connection.Close()

	var (
		users []*model.User
		err   error
	)
	err = connection.Order("id desc").Find(&users).Error
	return users, err
}

func (repository *userRepository) FetchByID(id int) (*model.User, error) {
	connection := conf.NewDBConnection()
	defer connection.Close()

	u := &model.User{ID: id}
	err := connection.First(u).Error
	return u, err
}

func (repository *userRepository) Update(u *model.User) (*model.User, error) {
	connection := conf.NewDBConnection()
	defer connection.Close()

	err := connection.Model(u).Update(u).Error
	return u, err
}

func (repository *userRepository) Delete(id int) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	u := &model.User{ID: id}
	err := connection.Delete(u).Error
	return err
}
