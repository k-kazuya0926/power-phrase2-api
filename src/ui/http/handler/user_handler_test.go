package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock
type mockUserUseCase struct {
	mock.Mock
}

// ユーザー登録
func (usecase *mockUserUseCase) CreateUser(name, email, password, imageFilePath string) (userID int, token string, err error) {
	args := usecase.Called(name, email, password, imageFilePath)
	return args.Int(0), args.String(1), args.Error(2)
}

// ログイン
func (usecase *mockUserUseCase) Login(email, password string) (userID int, token string, err error) {
	args := usecase.Called(email, password)
	return args.Int(0), args.String(1), args.Error(2)
}

// ユーザー詳細取得
func (usecase *mockUserUseCase) GetUser(id int) (*model.User, error) {
	args := usecase.Called(id)
	user, ok := args.Get(0).(*model.User)
	if ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

// ユーザー更新
func (usecase *mockUserUseCase) UpdateUser(userID int, name, email, password, imageFilePath string) error {
	return usecase.Called(userID, name, email, password, imageFilePath).Error(0)
}

// ユーザー削除
func (usecase *mockUserUseCase) DeleteUser(id int) error {
	return usecase.Called(id).Error(0)
}

func makeUser(id int) *model.User {
	return &model.User{
		ID:            id,
		CreatedAt:     time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
		UpdatedAt:     time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
		Name:          fmt.Sprintf("testuser%d", id),
		Email:         fmt.Sprintf("testuser%d@example.com", id),
		Password:      fmt.Sprintf("testuser%d", id),
		ImageFilePath: fmt.Sprintf("images/%d.png", id),
	}
}

// ユーザー登録成功
func TestCreateUser_success(t *testing.T) {
	// 1. Setup
	id := 1
	user := makeUser(id)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/users", strings.NewReader(string(jsonBytes)), rec)

	usecase := mockUserUseCase{}
	usecase.On("CreateUser", user.Name, user.Email, user.Password, user.ImageFilePath).Return(id, "token", nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.CreateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

// ユーザー登録バリデーションエラー
func TestCreateUser_error_validationError(t *testing.T) {
	cases := []struct {
		label    string
		name     string
		email    string
		password string
	}{
		{"name空", "", "testuser@example.com", "testuser"},
		{"email空", "testuser", "", "testuser"},
		{"email形式", "testuser", "testuserexample.com", "testuser"},
		{"password形式", "testuser", "testuser@example.com", ""},
	}

	for _, test := range cases {
		// 1. Setup
		user := makeUser(1)
		user.Name = test.name
		user.Email = test.email
		user.Password = test.password
		jsonBytes, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		c := createContext(echo.POST, "/users", strings.NewReader(string(jsonBytes)), rec)

		usecase := mockUserUseCase{}
		handler := NewUserHandler(&usecase)

		// 2. Exercise
		err = handler.CreateUser(c)

		// 3. Verify
		assert.NoError(t, err, test.label)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)

		// 4. Teardown
	}
}

// ユーザー登録ユースケースエラー
func TestCreateUser_error_usecaseError(t *testing.T) {
	// 1. Setup
	user := makeUser(1)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/users", strings.NewReader(string(jsonBytes)), rec)

	usecase := mockUserUseCase{}
	usecase.On("CreateUser", user.Name, user.Email, user.Password, user.ImageFilePath).Return(0, "", errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.CreateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// ログイン成功
func TestLogin_success(t *testing.T) {
	// 1. Setup
	email := "testuser@example.com"
	password := "testuser"
	reader := strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/users", reader, rec)

	usecase := mockUserUseCase{}
	usecase.On("Login", email, password).Return(1, "token", nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err := handler.Login(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEqual(t, "", rec.Body.String())

	// 4. Teardown
}

// ログインバリデーションエラー
func TestLogin_error_validationError(t *testing.T) {
	cases := []struct {
		label    string
		email    string
		password string
		message  string
	}{
		{"email空", "", "testuser", "\"Email：必須です。\"\n"},
		{"email形式", "testuserexample.com", "testuser", "\"Email：正しい形式で入力してください。\"\n"},
		{"password空", "testuser@example.com", "", "\"Password：必須です。\"\n"},
	}

	for _, test := range cases {
		// 1. Setup
		reader := strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, test.email, test.password))
		rec := httptest.NewRecorder()
		c := createContext(echo.POST, "/users", reader, rec)

		usecase := mockUserUseCase{}
		handler := NewUserHandler(&usecase)

		// 2. Exercise
		err := handler.Login(c)

		// 3. Verify
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, test.message, rec.Body.String())

		// 4. Teardown
	}
}

// ログインユースケースエラー
func TestLogin_error_usecaseError(t *testing.T) {
	// 1. Setup
	email := "testuser@example.com"
	password := "testuser"
	reader := strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/users", reader, rec)

	usecase := mockUserUseCase{}
	usecase.On("Login", email, password).Return(0, "", errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err := handler.Login(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// ユーザー詳細取得成功
func TestGetUser_success(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.GET, "/users", nil, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	expectedUser := makeUser(id)

	usecase := mockUserUseCase{}
	usecase.On("GetUser", id).Return(expectedUser, nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err := handler.GetUser(c)

	// 3. Verify
	assert.NoError(t, err)
	user := &model.User{}
	json.Unmarshal(rec.Body.Bytes(), user)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.ImageFilePath, user.ImageFilePath)

	// 4. Teardown
}

// ユーザー詳細取得バリデーションエラー
func TestGetUser_error_validationError(t *testing.T) {
	cases := []struct {
		label   string
		id      interface{}
		message string
	}{
		{"必須", "", "\"ID：数値で入力してください。\"\n"},
		{"型", "a", "\"ID：数値で入力してください。\"\n"},
		{"下限", 0, "\"ID：1以上の値を入力してください。\"\n"},
	}

	for _, test := range cases {
		// 1. Setup
		rec := httptest.NewRecorder()
		c := createContext(echo.GET, "/users", nil, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		usecase := mockUserUseCase{}
		handler := NewUserHandler(&usecase)

		// 2. Exercise
		err := handler.GetUser(c)

		// 3. Verify
		assert.NoError(t, err, test.label)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
		assert.Equal(t, test.message, rec.Body.String(), test.label)

		// 4. Teardown
	}
}

// ユーザー詳細取得ユースケースエラー
func TestGetUser_error_usecaseError(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.GET, "/users", nil, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	usecase := mockUserUseCase{}
	usecase.On("GetUser", id).Return(nil, errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err := handler.GetUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// ユーザー更新成功
func TestUpdateUser_success(t *testing.T) {
	// 1. Setup
	user := makeUser(1)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	c := createContext(echo.PUT, "/users", strings.NewReader(string(jsonBytes)), rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(1))

	usecase := mockUserUseCase{}
	usecase.On("UpdateUser", user.ID, user.Name, user.Email, user.Password, user.ImageFilePath).Return(nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.UpdateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

// ユーザー更新バリデーションエラー
func TestUpdateUser_error_validationError(t *testing.T) {
	cases := []struct {
		label   string
		id      interface{}
		name    string
		email   string
		message string
	}{
		{"ID空", "", "testuser", "testuser@example.com", "\"ID：数値で入力してください。\"\n"},
		{"ID形式", "a", "testuser", "testuser@example.com", "\"ID：数値で入力してください。\"\n"},
		{"ID下限", 0, "testuser", "testuser@example.com", "\"ID：必須です。\"\n"},
		{"Name空", 1, "", "testuser@example.com", "\"Name：必須です。\"\n"},
		{"Email空", 1, "testuser", "", "\"Email：必須です。\"\n"},
		{"Email形式", 1, "testuser", "testuserexample.com", "\"Email：正しい形式で入力してください。\"\n"},
	}

	for _, test := range cases {
		// 1. Setup
		rec := httptest.NewRecorder()
		c := createContext(echo.PUT, "/users", strings.NewReader(fmt.Sprintf(`{
			"name": "%s",
			"email": "%s",
			"password": "testuser"
		}`, test.name, test.email)), rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		usecase := mockUserUseCase{}
		handler := NewUserHandler(&usecase)

		// 2. Exercise
		err := handler.UpdateUser(c)

		// 3. Verify
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
		assert.Equal(t, test.message, rec.Body.String(), test.label)

		// 4. Teardown
	}
}

// ユーザー更新ユースケースエラー
func TestUpdateUser_error_usecaseError(t *testing.T) {
	// 1. Setup
	user := makeUser(1)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	c := createContext(echo.PUT, "/users", strings.NewReader(string(jsonBytes)), rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	usecase := mockUserUseCase{}
	usecase.On("UpdateUser", user.ID, user.Name, user.Email, user.Password, user.ImageFilePath).Return(errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.UpdateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// ユーザー削除成功
func TestDeleteUser_success(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.DELETE, "/users", nil, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	usecase := mockUserUseCase{}
	usecase.On("DeleteUser", id).Return(nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err := handler.DeleteUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

// ユーザー削除バリデーションエラー
func TestDeleteUser_error_validationError(t *testing.T) {
	cases := []struct {
		label   string
		id      interface{}
		message string
	}{
		{"ID空", "", "\"ID：数値で入力してください。\"\n"},
		{"ID形式", "a", "\"ID：数値で入力してください。\"\n"},
		{"ID下限", 0, "\"ID：1以上の値を入力してください。\"\n"},
	}

	for _, test := range cases {
		// 1. Setup
		rec := httptest.NewRecorder()
		c := createContext(echo.DELETE, "/users", nil, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		usecase := mockUserUseCase{}
		handler := NewUserHandler(&usecase)

		// 2. Exercise
		err := handler.DeleteUser(c)

		// 3. Verify
		assert.NoError(t, err, test.label)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
		assert.Equal(t, test.message, rec.Body.String(), test.label)

		// 4. Teardown
	}
}

// ユーザー削除ユースケースエラー
func TestDeleteUser_error_usecaseError(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.DELETE, "/users", nil, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	usecase := mockUserUseCase{}
	usecase.On("DeleteUser", id).Return(errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err := handler.DeleteUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}
