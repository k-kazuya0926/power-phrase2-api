package datastore

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

type userRepository struct {
	conn *gorm.DB
}

// NewUserRepository UserRepositoryを取得します.
func NewUserRepository(conn *gorm.DB) repository.UserRepository {
	return &userRepository{conn}
}

func (repository *userRepository) Fetch(ctx context.Context) ([]*model.User, error) {
	var (
		users []*model.User
		err   error
	)
	err = repository.conn.Order("id desc").Find(&users).Error
	return users, err
}

func (repository *userRepository) FetchByID(ctx context.Context, id int) (*model.User, error) {
	u := &model.User{ID: id}
	err := repository.conn.First(u).Error
	return u, err
}

func (repository *userRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	err := repository.conn.Create(u).Error
	return u, err
}

func (repository *userRepository) Update(ctx context.Context, u *model.User) (*model.User, error) {
	err := repository.conn.Model(u).Update(u).Error
	return u, err
}

func (repository *userRepository) Delete(ctx context.Context, id int) error {
	u := &model.User{ID: id}
	err := repository.conn.Delete(u).Error
	return err
}
