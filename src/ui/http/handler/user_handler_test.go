package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/validator"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock
type mockUserUseCase struct {
	mock.Mock
}

func (usecase *mockUserUseCase) CreateUser(name, email, password, imageURL string) (err error) {
	return usecase.Called(name, email, password, imageURL).Error(0)
}

func (usecase *mockUserUseCase) Login(email, password string) (userID int, token string, err error) {
	args := usecase.Called(email, password)
	return args.Int(0), args.String(1), args.Error(2)
}

func (usecase *mockUserUseCase) GetUser(id int) (*model.User, error) {
	args := usecase.Called(id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*model.User), args.Error(1)
}

func (usecase *mockUserUseCase) UpdateUser(userID int, name, email, password, imageURL string) error {
	return usecase.Called(userID, name, email, password, imageURL).Error(0)
}

func (usecase *mockUserUseCase) DeleteUser(id int) error {
	return usecase.Called(id).Error(0)
}

func getMockUser(id int) *model.User {
	user := &model.User{
		ID:        id,
		CreatedAt: time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
		UpdatedAt: time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
		Name:      fmt.Sprintf("testuser%d", id),
		Email:     fmt.Sprintf("testuser%d@example.com", id),
		Password:  fmt.Sprintf("testuser%d", id),
		ImageURL:  fmt.Sprintf("http://www.example.com/%d", id),
	}
	return user
}

// ユーザー登録テスト
func TestCreateUser_success(t *testing.T) {
	// 1. Setup
	user := getMockUser(1)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/users", strings.NewReader(string(jsonBytes)), rec)

	usecase := mockUserUseCase{}
	usecase.On("CreateUser", user.Name, user.Email, user.Password, user.ImageURL).Return(nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.CreateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

func createContext(method, path string, body io.Reader, rec *httptest.ResponseRecorder) echo.Context {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e := echo.New()
	e.Validator = validator.NewValidator()
	return e.NewContext(req, rec)
}

func TestCreateUser_error_emptyName(t *testing.T) {
	// 1. Setup
	user := getMockUser(1)
	user.Name = ""
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
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// 4. Teardown
}

func TestCreateUser_error_invalidEmail(t *testing.T) {
	cases := []struct {
		label string
		email string
	}{
		{"空", ""},
		{"形式", "testuserexample.com"},
	}

	for _, test := range cases {
		// 1. Setup
		user := getMockUser(1)
		user.Email = test.email
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

func TestCreateUser_error_emptyPassword(t *testing.T) {
	// 1. Setup
	user := getMockUser(1)
	user.Password = ""
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
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// 4. Teardown
}

func TestCreateUser_error_usecaseError(t *testing.T) {
	// 1. Setup
	user := getMockUser(1)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/users", strings.NewReader(string(jsonBytes)), rec)

	usecase := mockUserUseCase{}
	usecase.On("CreateUser", user.Name, user.Email, user.Password, user.ImageURL).Return(errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.CreateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// ログインテスト
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

// ユーザー詳細テスト
func TestGetUser_success(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.GET, "/users", nil, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	expectedUser := getMockUser(id)

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
	assert.Equal(t, expectedUser, user)

	// 4. Teardown
}

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

func TestUpdateUser_success(t *testing.T) {
	// 1. Setup
	user := getMockUser(1)
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
	usecase.On("UpdateUser", user.ID, user.Name, user.Email, user.Password, user.ImageURL).Return(nil)
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.UpdateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

func TestUpdateUser_error_invalidID(t *testing.T) {
	cases := []struct {
		label   string
		id      interface{}
		message string
	}{
		{"必須", "", "\"ID：数値で入力してください。\"\n"},
		{"型", "a", "\"ID：数値で入力してください。\"\n"},
		{"下限", 0, "\"ID：必須です。\"\n"},
	}

	for _, test := range cases {
		// 1. Setup
		rec := httptest.NewRecorder()
		c := createContext(echo.PUT, "/users", strings.NewReader(`{}`), rec)
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

func TestUpdateUser_error_usecaseError(t *testing.T) {
	// 1. Setup
	user := getMockUser(1)
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
	usecase.On("UpdateUser", user.ID, user.Name, user.Email, user.Password, user.ImageURL).Return(errors.New("error"))
	handler := NewUserHandler(&usecase)

	// 2. Exercise
	err = handler.UpdateUser(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// ユーザー削除テスト
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

func TestDeleteUser_error_invalidID(t *testing.T) {
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
