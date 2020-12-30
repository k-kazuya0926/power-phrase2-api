package datastore

import (
	"fmt"
	"testing"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
)

// ユーザー登録
func TestUserRepository_Create(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := makeUserForInput(1)

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
	assert.Equal(t, userForInput.ImageFilePath, user.ImageFilePath)

	// 4. Teardown
	teardown(db)
}

// メールアドレスによるユーザー詳細取得
func TestUserRepository_FetchByEmail(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	userForInput2 := makeUserForInput(2)
	db.Create(&userForInput2)

	// 2. Exercise
	user, err := repository.FetchByEmail(userForInput.Email)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, userForInput.Name, user.Name)
	assert.Equal(t, userForInput.Email, user.Email)
	assert.Equal(t, userForInput.Password, user.Password)
	assert.Equal(t, userForInput.ImageFilePath, user.ImageFilePath)

	// 4. Teardown
	teardown(db)
}

// IDによるユーザー詳細取得
func TestUserRepository_FetchById(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	repository := &userRepository{}

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
	assert.Equal(t, userForInput.ImageFilePath, actualUser.ImageFilePath)

	// 4. Teardown
	teardown(db)
}

// ユーザー更新
func TestUserRepository_Update(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)
	userForInput.Name = "testuser2"
	userForInput.Email = "testuser2@example.com"
	userForInput.Password = "testuser2"
	userForInput.ImageFilePath = "images/2.png"

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
	assert.Equal(t, userForInput.ImageFilePath, user.ImageFilePath)

	// 4. Teardown
	teardown(db)
}

// ユーザー削除
func TestUserRepository_Delete(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	repository := &userRepository{}
	userForInput := makeUserForInput(1)
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

// 入力用ユーザー生成
func makeUserForInput(id int) *model.User {
	user := &model.User{
		Name:          fmt.Sprintf("testuser%d", id),
		Email:         fmt.Sprintf("testuser%d@example.com", id),
		Password:      fmt.Sprintf("testuser%d", id),
		ImageFilePath: fmt.Sprintf("images/%d.png", id),
	}
	return user
}
