package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

var handler UserHandler
var e *echo.Echo

func setUp() {
	// set stub
	usecase := &mockUserUseCase{}
	handler = NewUserHandler(usecase)

	e = echo.New()
	e.Validator = validator.NewValidator()
}

func tearDown() {
}

type mockUserUseCase struct{}

func (u *mockUserUseCase) CreateUser(name, email, password, imageURL string) (err error) {
	return nil
}

func (u *mockUserUseCase) Login(email, password string) (token string, err error) {
	return token, err
}

func (u *mockUserUseCase) GetUser(id int) (*model.User, error) {
	return getMockUser(id), nil
}

func (u *mockUserUseCase) UpdateUser(userID int, name, email, password, imageURL string) error {
	return nil
}

func (u *mockUserUseCase) DeleteUser(id int) error {
	return nil
}

// TODO 見直し
// func getMockUsers(n int) []*model.User {
// 	ret := []*model.User{}
// 	for i := 0; i < n; i++ {
// 		u := getMockUser(int(i))
// 		ret = append(ret, u)
// 	}
// 	return ret
// }

func getMockUser(id int) *model.User {
	u := &model.User{
		ID:        id,
		Name:      fmt.Sprintf("testuser%d", id),
		Email:     fmt.Sprintf("testuser%d@example.com", id),
		Password:  fmt.Sprintf("testuser%d", id),
		ImageURL:  fmt.Sprintf("http://www.example.com/%d", id),
		CreatedAt: time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
		UpdatedAt: time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
	}
	return u
}

func getMockUserNoID() *model.User {
	u := &model.User{
		Name:      fmt.Sprintf("name_%d", 1),
		CreatedAt: time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
		UpdatedAt: time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local),
	}
	return u
}

func TestGetUser_success(t *testing.T) {
	cases := []struct {
		ID   int
		User *model.User
	}{
		{1, getMockUser(1)},
	}

	for _, test := range cases {
		req := httptest.NewRequest(echo.GET, "/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.ID))

		// assertions
		if assert.NoError(t, handler.GetUser(c)) {
			user := &model.User{}
			if err := json.Unmarshal(rec.Body.Bytes(), user); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, test.User, user)
		}
	}
}

func TestGetUser_error(t *testing.T) {
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

// TODO
// func TestUserHandler_CreateUser(t *testing.T) {
// 	// expected
// 	expected := getMockUser(1)

// 	// set stub
// 	usecase := &mockUserUseCase{}
// 	handler := NewUserHandler(usecase)

// 	e := echo.New()
// 	jsonBytes, err := json.Marshal(getMockUserNoID())
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(string(jsonBytes)))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	// Assertions
// 	if assert.NoError(t, handler.CreateUser(c)) {
// 		user := &model.User{}
// 		if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
// 			t.Fatal(err)
// 		}

// 		t.Log(rec.Code)
// 		assert.Equal(t, http.StatusCreated, rec.Code)
// 		t.Log(user)
// 		assert.Equal(t, expected, user)
// 	}
// }

// TODO
// type updateUserTest struct {
// 	ID       int
// 	UserName string
// }

// var updateUserTests = []updateUserTest{
// 	{math.MaxInt8, fmt.Sprintf("name_%d_updated", math.MaxInt8)},
// 	{math.MaxInt16, fmt.Sprintf("name_%d_updated", math.MaxInt16)},
// 	{math.MaxInt32, fmt.Sprintf("name_%d_updated", math.MaxInt32)},
// 	{math.MaxInt64, fmt.Sprintf("name_%d_updated", math.MaxInt64)},
// }

// func TestUserHandler_UpdateUser(t *testing.T) {
// 	// set stub
// 	usecase := &mockUserUseCase{}
// 	handler := NewUserHandler(usecase)

// 	for _, test := range updateUserTests {
// 		// set request
// 		e := echo.New()
// 		req := httptest.NewRequest(echo.PUT, "/users", nil)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)
// 		c.SetPath("/users/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues(fmt.Sprint(test.ID))

// 		// assertions
// 		if assert.NoError(t, handler.UpdateUser(c)) {
// 			user := &model.User{}
// 			if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
// 				t.Fatal(err)
// 			}
// 			t.Log(rec.Code)
// 			assert.Equal(t, http.StatusOK, rec.Code)
// 			t.Log(user)
// 			assert.Equal(t, test.UserName, user.Name)
// 		}
// 	}
// }

// TODO
// func TestUserHandler_DeleteUser(t *testing.T) {
// 	// set stub
// 	usecase := &mockUserUseCase{}
// 	handler := NewUserHandler(usecase)

// 	// set request
// 	e := echo.New()
// 	req := httptest.NewRequest(echo.DELETE, "/users", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	c.SetPath("/users/:id")
// 	c.SetParamNames("id")
// 	c.SetParamValues(fmt.Sprint(1))

// 	// assertions
// 	if assert.NoError(t, handler.DeleteUser(c)) {
// 		t.Log(rec.Code)
// 		assert.Equal(t, http.StatusNoContent, rec.Code)
// 	}
// }
