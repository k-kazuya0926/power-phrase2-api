package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock
type mockUserRepository struct {
	mock.Mock
}

func (repository *mockUserRepository) Create(user *model.User) error {
	args := repository.Called(user)
	user.ID = 1
	return args.Error(0)
}

func (repository *mockUserRepository) FetchByEmail(email string) (*model.User, error) {
	args := repository.Called(email)
	user, ok := args.Get(0).(*model.User)
	if ok {
		return user, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (repository *mockUserRepository) FetchByID(id int) (*model.User, error) {
	args := repository.Called(id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (repository *mockUserRepository) Update(user *model.User) error {
	return repository.Called(user).Error(0)
}

func (repository *mockUserRepository) Delete(id int) error {
	return repository.Called(id).Error(0)
}

// 入力用ユーザー
func getMockUserForInput(id int) *model.User {
	user := &model.User{
		Name:          fmt.Sprintf("testuser%d", id),
		Email:         fmt.Sprintf("testuser%d@example.com", id),
		Password:      fmt.Sprintf("testuser%d", id),
		ImageFilePath: fmt.Sprintf("/images/%d.png", id),
	}
	return user
}

// DBから取得されたユーザー
func getMockUserForRead(id int) *model.User {
	user := getMockUserForInput(id)
	user.ID = id
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordHash)
	user.CreatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	user.UpdatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	return user
}

// ユーザー登録テスト
func TestCreateUser_success(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	user := getMockUserForInput(id)
	repository.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	// 2. Exercise
	userID, token, err := usecase.CreateUser(user.Name, user.Email, user.Password, user.ImageFilePath)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, id, userID)
	assert.NotEqual(t, "", token)

	// 4. Teardown
}

func TestCreateUser_error(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	user := getMockUserForInput(id)
	repository.On("Create", mock.AnythingOfType("*model.User")).Return(errors.New("error"))

	// 2. Exercise
	userID, token, err := usecase.CreateUser(user.Name, user.Email, user.Password, user.ImageFilePath)

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Equal(t, "", token)

	// 4. Teardown
}

// ログインテスト
func TestLogin_success(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	userForInput := getMockUserForInput(id)
	userForRead := getMockUserForRead(id)
	repository.On("FetchByEmail", userForInput.Email).Return(userForRead, nil)

	// 2. Exercise
	userID, token, err := usecase.Login(userForInput.Email, userForInput.Password)

	// 3. Verify
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userID)
	assert.NotEqual(t, "", token)

	// 4. Teardown
}

func TestLogin_error_invalidPassword(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	userForInput := getMockUserForInput(id)
	userForRead := getMockUserForRead(id)
	repository.On("FetchByEmail", userForInput.Email).Return(userForRead, nil)

	// 2. Exercise
	userID, token, err := usecase.Login(userForInput.Email, "invalid")

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Equal(t, "", token)

	// 4. Teardown
}

func TestLogin_error_repositoryError(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	userForInput := getMockUserForInput(id)
	repository.On("FetchByEmail", userForInput.Email).Return(nil, errors.New("error"))

	// 2. Exercise
	userID, token, err := usecase.Login(userForInput.Email, userForInput.Password)

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Equal(t, "", token)

	// 4. Teardown
}

// ユーザー詳細テスト
func TestGetUser_success(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	expected := getMockUserForRead(id)
	expected.Password = ""
	repository.On("FetchByID", id).Return(expected, nil)

	// 2. Exercise
	user, err := usecase.GetUser(id)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.CreatedAt, user.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, user.UpdatedAt)
	assert.Equal(t, expected.DeletedAt, user.DeletedAt)
	assert.Equal(t, expected.Name, user.Name)
	assert.Equal(t, expected.Password, user.Password)
	assert.Equal(t, expected.Email, user.Email)
	assert.Equal(t, expected.ImageFilePath, user.ImageFilePath)

	// 4. Teardown
}

func TestGetUser_error(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	repository.On("FetchByID", id).Return(nil, errors.New("error"))

	// 2. Execise
	user, err := usecase.GetUser(id)

	// 3. Verify
	assert.Error(t, err)
	assert.Empty(t, user)

	// 4. Teardown
}

// ユーザー更新テスト
func TestUpdateUser_success(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	user := getMockUserForInput(id)
	repository.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	// 2. Exercise
	err := usecase.UpdateUser(id, user.Name, user.Email, user.Password, user.ImageFilePath)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestUpdateUser_error(t *testing.T) {
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	user := getMockUserForInput(id)
	repository.On("Update", mock.AnythingOfType("*model.User")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.UpdateUser(id, user.Name, user.Email, user.Password, user.ImageFilePath)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// ユーザー削除テスト
func TestDeleteUser_success(t *testing.T) {
	// 1. Setup
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(nil)

	// 2. Exercise
	err := usecase.DeleteUser(id)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestDeleteUser_error(t *testing.T) {
	repository := mockUserRepository{}
	usecase := NewUserUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.DeleteUser(id)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}
