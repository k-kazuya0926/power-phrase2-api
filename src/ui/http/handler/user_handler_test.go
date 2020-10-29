package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/validator"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

var uc mockUserUseCase
var handler UserHandler
var e *echo.Echo

func setUp() {
	// set stub
	uc = mockUserUseCase{}
	handler = NewUserHandler(&uc)

	e = echo.New()
	e.Validator = validator.NewValidator()
}

func tearDown() {
}

// Mock
type mockUserUseCase struct {
	returnsError bool
}

func (uc *mockUserUseCase) CreateUser(name, email, password, imageURL string) (err error) {
	if uc.returnsError {
		return errors.New("error")
	}

	return nil
}

func (uc *mockUserUseCase) Login(email, password string) (token string, err error) {
	if uc.returnsError {
		return "", errors.New("error")
	}

	return token, err
}

func (uc *mockUserUseCase) GetUser(id int) (*model.User, error) {
	if uc.returnsError {
		return nil, errors.New("error")
	}

	return getMockUser(id), nil
}

func (uc *mockUserUseCase) UpdateUser(userID int, name, email, password, imageURL string) error {
	if uc.returnsError {
		return errors.New("error")
	}

	return nil
}

func (uc *mockUserUseCase) DeleteUser(id int) error {
	if uc.returnsError {
		return errors.New("error")
	}

	return nil
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
	jsonBytes, err := json.Marshal(getMockUser(1))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, handler.CreateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateUser_error_emptyName(t *testing.T) {
	user := getMockUser(1)
	user.Name = ""
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// assertions
	if assert.NoError(t, handler.CreateUser(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
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
		user := getMockUser(1)
		user.Email = test.email
		jsonBytes, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(string(jsonBytes)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// assertions
		if assert.NoError(t, handler.CreateUser(c)) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestCreateUser_error_emptyPassword(t *testing.T) {
	user := getMockUser(1)
	user.Password = ""
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// assertions
	if assert.NoError(t, handler.CreateUser(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	}
}

func TestCreateUser_error_usecaseError(t *testing.T) {
	uc.returnsError = true

	user := getMockUser(1)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(string(jsonBytes)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// assertions
	if assert.NoError(t, handler.CreateUser(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	uc.returnsError = false
}

// ログインテスト
func TestLogin_success(t *testing.T) {
	reader := strings.NewReader(`{"email": "testuser@example.com", "password": "testuser"}`)
	req := httptest.NewRequest(echo.POST, "/users", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotEqual(t, "", rec.Body.String())
	}
}

func TestLogin_error_invalidEmail(t *testing.T) {
	cases := []struct {
		label string
		email string
	}{
		{"空", ""},
		{"形式", "testuserexample.com"},
	}

	for _, test := range cases {
		reader := strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "testuser"}`, test.email))
		req := httptest.NewRequest(echo.POST, "/users", reader)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assertions
		if assert.NoError(t, handler.Login(c)) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	}
}

func TestLogin_error_emptyPassword(t *testing.T) {
	reader := strings.NewReader(`{"email": "testuser@exampl.com", "password": ""}`)
	req := httptest.NewRequest(echo.POST, "/users", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, "\"Password：必須です。\"\n", rec.Body.String())
	}
}

func TestLogin_error_usecaseError(t *testing.T) {
	uc.returnsError = true
	reader := strings.NewReader(`{"email": "testuser@example.com", "password": "testuser"}`)
	req := httptest.NewRequest(echo.POST, "/users", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	uc.returnsError = false
}

// ユーザー詳細テスト
func TestGetUser_success(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	id := 1
	c.SetParamValues(fmt.Sprint(id))

	// assertions
	if assert.NoError(t, handler.GetUser(c)) {
		user := &model.User{}
		if err := json.Unmarshal(rec.Body.Bytes(), user); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, getMockUser(id), user)
	}
}

func TestGetUser_error_invalidID(t *testing.T) {
	cases := []struct {
		label   string
		id      interface{}
		message string
	}{
		{"必須", "", "\"ID：数値で入力してください。\"\n"},
		{"型", "a", "\"ID：数値で入力してください。\"\n"},
		{"下限", 0, "\"ID：1以上の値を入力してください。\"\n"},
	}

	req := httptest.NewRequest(echo.GET, "/users", nil)
	for _, test := range cases {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		// assertions
		if assert.NoError(t, handler.GetUser(c), test.label) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
			assert.Equal(t, test.message, rec.Body.String(), test.label)
		}
	}
}

func TestGetUser_error_usecaseError(t *testing.T) {
	uc.returnsError = true
	req := httptest.NewRequest(echo.GET, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(1))

	// assertions
	if assert.NoError(t, handler.GetUser(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	uc.returnsError = false
}

func TestUpdateUser_success(t *testing.T) {
	reader := strings.NewReader(`{}`)
	req := httptest.NewRequest(echo.PUT, "/users", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(1))

	// assertions
	if assert.NoError(t, handler.UpdateUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
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

	reader := strings.NewReader(`{}`)
	req := httptest.NewRequest(echo.PUT, "/users", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for _, test := range cases {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		// assertions
		if assert.NoError(t, handler.UpdateUser(c)) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
			assert.Equal(t, test.message, rec.Body.String(), test.label)
		}
	}
}

func TestUpdateUser_error_usecaseError(t *testing.T) {
	uc.returnsError = true
	reader := strings.NewReader(`{}`)
	req := httptest.NewRequest(echo.PUT, "/users", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(1))

	// assertions
	if assert.NoError(t, handler.UpdateUser(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	uc.returnsError = false
}

// ユーザー削除テスト
func TestDeleteUser_success(t *testing.T) {
	req := httptest.NewRequest(echo.DELETE, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(1))

	// assertions
	if assert.NoError(t, handler.DeleteUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
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

	req := httptest.NewRequest(echo.DELETE, "/users", nil)
	for _, test := range cases {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		// assertions
		if assert.NoError(t, handler.DeleteUser(c), test.label) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
			assert.Equal(t, test.message, rec.Body.String(), test.label)
		}
	}
}

func TestDeleteUser_error_usecaseError(t *testing.T) {
	uc.returnsError = true
	req := httptest.NewRequest(echo.DELETE, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(1))

	// assertions
	if assert.NoError(t, handler.DeleteUser(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	}

	uc.returnsError = false
}
