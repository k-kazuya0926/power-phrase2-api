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

func (repository *userRepository) Create(user *model.User) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	return connection.Create(user).Error
}

func (repository *userRepository) FetchByEmail(email string) (*model.User, error) {
	var user model.User

	connection := conf.NewDBConnection()
	defer connection.Close()

	err := connection.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (repository *userRepository) FetchByID(id int) (*model.User, error) {
	connection := conf.NewDBConnection()
	defer connection.Close()

	u := model.User{ID: id}
	err := connection.First(&u).Error
	u.Password = ""

	return &u, err
}

func (repository *userRepository) Update(u *model.User) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	return connection.Model(u).Update(u).Error
}

func (repository *userRepository) Delete(id int) error {
	connection := conf.NewDBConnection()
	defer connection.Close()

	user := model.User{ID: id}
	return connection.Delete(&user).Error
}
