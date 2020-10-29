package datastore

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
)

func teardown(db *gorm.DB) {
	db.DropTable(&model.Post{})
	db.DropTable(&model.User{})
}

func TestCreate(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := getMockUserForInput(1)

	// 2. Exercise
	err := repository.Create(userForInput)

	// 3. Verify
	assert.NoError(t, err)

	// 件数
	var count int
	db.Table("users").Count(&count)
	assert.Equal(t, 1, count)

	// 内容
	user := model.User{}
	db.First(&user)
	assert.Equal(t, userForInput.Name, user.Name)
	assert.Equal(t, userForInput.Email, user.Email)
	assert.Equal(t, userForInput.Password, user.Password)
	assert.Equal(t, userForInput.ImageURL, user.ImageURL)

	// 4. Teardown
	teardown(db)
}

func TestFetchByEmail(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	userForInput2 := getMockUserForInput(2)
	db.Create(&userForInput2)

	// 2. Exercise
	user, err := repository.FetchByEmail(userForInput.Email)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, userForInput.Name, user.Name)
	assert.Equal(t, userForInput.Email, user.Email)
	assert.Equal(t, userForInput.Password, user.Password)
	assert.Equal(t, userForInput.ImageURL, user.ImageURL)

	// 4. Teardown
	teardown(db)
}

func TestFetchById(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	// 2. Exercise
	actualUser, err := repository.FetchByID(userForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, userForInput.ID, actualUser.ID)
	assert.Equal(t, userForInput.Name, actualUser.Name)
	assert.Equal(t, userForInput.Name, actualUser.Name)
	assert.Equal(t, userForInput.Email, actualUser.Email)
	assert.Equal(t, "", actualUser.Password)
	assert.Equal(t, userForInput.ImageURL, actualUser.ImageURL)

	// 4. Teardown
	teardown(db)
}

func TestUpdate(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)
	userForInput.Name = "testuser2"
	userForInput.Email = "testuser2@example.com"
	userForInput.Password = "testuser2"
	userForInput.ImageURL = "http://www.example.com/2"
	db.Update(userForInput)

	// 2. Exercise
	err := repository.Update(userForInput)

	// 3. Verify
	assert.NoError(t, err)

	// 件数
	var count int
	db.Table("users").Count(&count)
	assert.Equal(t, 1, count)

	// 内容
	user := model.User{}
	db.First(&user)
	assert.Equal(t, userForInput.Name, user.Name)
	assert.Equal(t, userForInput.Email, user.Email)
	assert.Equal(t, userForInput.Password, user.Password)
	assert.Equal(t, userForInput.ImageURL, user.ImageURL)

	// 4. Teardown
	teardown(db)
}

func TestDelete(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	// 2. Exercise
	err := repository.Delete(userForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	assert.Error(t, db.First(&userForInput).Error)

	// 4. Teardown
	teardown(db)
}

// 入力用ユーザー
func getMockUserForInput(id int) *model.User {
	user := &model.User{
		Name:     fmt.Sprintf("testuser%d", id),
		Email:    fmt.Sprintf("testuser%d@example.com", id),
		Password: fmt.Sprintf("testuser%d", id),
		ImageURL: fmt.Sprintf("http://www.example.com/%d", id),
	}
	return user
}
