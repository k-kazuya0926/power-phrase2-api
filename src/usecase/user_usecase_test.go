package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Mock
type mockUserRepository struct {
	returnsError bool
}

func (ur *mockUserRepository) Create(user *model.User) error {
	if ur.returnsError {
		return errors.New("error")
	}

	return nil
}

func (ur *mockUserRepository) FetchByEmail(email string) (*model.User, error) {
	if ur.returnsError {
		return nil, errors.New("error")
	}

	return getMockUserForRead(1), nil
}

func (ur *mockUserRepository) FetchByID(id int) (*model.User, error) {
	if ur.returnsError {
		return nil, errors.New("error")
	}

	return getMockUserForRead(id), nil
}

func (ur *mockUserRepository) Update(user *model.User) error {
	if ur.returnsError {
		return errors.New("error")
	}

	return nil
}

func (ur *mockUserRepository) Delete(id int) error {
	if ur.returnsError {
		return errors.New("error")
	}

	return nil
}

var ur = mockUserRepository{}
var uc = NewUserUseCase(&ur)

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

// DBから取得されたユーザー
func getMockUserForRead(id int) *model.User {
	user := getMockUserForInput(id)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordHash)
	user.CreatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	user.UpdatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	return user
}

// ユーザー登録テスト
func TestCreateUser_success(t *testing.T) {
	user := getMockUserForInput(1)

	// Assertions
	assert.NoError(t, uc.CreateUser(user.Name, user.Email, user.Password, user.ImageURL))
}

func TestCreateUser_error(t *testing.T) {
	ur.returnsError = true
	user := getMockUserForInput(1)

	// Assertions
	assert.Error(t, uc.CreateUser(user.Name, user.Email, user.Password, user.ImageURL))

	ur.returnsError = false
}

// ログインテスト
func TestLogin_success(t *testing.T) {
	user := getMockUserForInput(1)
	token, err := uc.Login(user.Email, user.Password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEqual(t, "", token)
}

func TestLogin_error_invalidPassword(t *testing.T) {
	user := getMockUserForInput(1)
	token, err := uc.Login(user.Email, "invalid")

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestLogin_error_repositoryError(t *testing.T) {
	ur.returnsError = true

	user := getMockUserForInput(1)
	token, err := uc.Login(user.Email, user.Password)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "", token)

	ur.returnsError = false
}

// ユーザー詳細テスト
func TestGetUser_success(t *testing.T) {
	id := 1
	expected := getMockUserForRead(id)

	user, err := uc.GetUser(id)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.Name, user.Name)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.CreatedAt, user.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, user.UpdatedAt)
	assert.Equal(t, expected.DeletedAt, user.DeletedAt)
}

func TestGetUser_error(t *testing.T) {
	ur.returnsError = true
	id := 1
	user, err := uc.GetUser(id)

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, user)

	ur.returnsError = false
}

// ユーザー更新テスト
func TestUpdateUser_success(t *testing.T) {
	id := 1
	user := getMockUserForInput(id)

	// Assertions
	assert.NoError(t, uc.UpdateUser(id, user.Name, user.Email, user.Password, user.ImageURL))
}

func TestUpdateUser_error(t *testing.T) {
	ur.returnsError = true
	id := 1
	user := getMockUserForInput(id)

	// Assertions
	assert.Error(t, uc.UpdateUser(id, user.Name, user.Email, user.Password, user.ImageURL))

	ur.returnsError = false
}

// ユーザー削除テスト
func TestDeleteUser_success(t *testing.T) {
	id := 1

	// Assertions
	assert.NoError(t, uc.DeleteUser(id))
}

func TestDeleteUser_error(t *testing.T) {
	ur.returnsError = true
	id := 1

	// Assertions
	assert.Error(t, uc.DeleteUser(id))

	ur.returnsError = false
}
