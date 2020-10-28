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
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(user).Error
}

func (repository *userRepository) FetchByEmail(email string) (*model.User, error) {
	var user model.User

	db := conf.NewDBConnection()
	defer db.Close()

	err := db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (repository *userRepository) FetchByID(id int) (*model.User, error) {
	db := conf.NewDBConnection()
	defer db.Close()

	u := model.User{ID: id}
	if err := db.First(&u).Error; err != nil {
		return nil, err
	}
	u.Password = ""

	return &u, nil
}

func (repository *userRepository) Update(u *model.User) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Model(u).Update(u).Error
}

func (repository *userRepository) Delete(id int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	user := model.User{ID: id}
	return db.Delete(&user).Error
}
