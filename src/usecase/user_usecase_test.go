package usecase

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

var ur mockUserRepository
var uc UserUseCase

func setUp() {
	ur = mockUserRepository{}
	uc = NewUserUseCase(&ur)
}

func tearDown() {
}

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
func TestCreate_success(t *testing.T) {
	user := getMockUserForInput(1)

	// Assertions
	assert.NoError(t, uc.CreateUser(user.Name, user.Email, user.Password, user.ImageURL))
}

func TestCreate_error(t *testing.T) {
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

// TODO GetUser

// TODO UpdateUser

// TODO DeleteUser
