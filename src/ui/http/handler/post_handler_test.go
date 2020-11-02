package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock
type mockPostUseCase struct {
	mock.Mock
}

func (usecase *mockPostUseCase) CreatePost(userID int, title, speaker, detail, movieURL string) (err error) {
	return usecase.Called(userID, title, speaker, detail, movieURL).Error(0)
}

func (usecase *mockPostUseCase) GetPosts(limit, offset int, keyword string) (totalCount int, posts []*model.GetPostResult, err error) {
	args := usecase.Called(limit, offset, keyword)
	return args.Int(0), args.Get(1).([]*model.GetPostResult), args.Error(2)
}

func (usecase *mockPostUseCase) GetPost(id int) (*model.Post, error) {
	args := usecase.Called(id)
	post := args.Get(0)
	if post == nil {
		return nil, args.Error(1)
	}
	return post.(*model.Post), args.Error(1)
}

func (usecase *mockPostUseCase) UpdatePost(ID int, title, speaker, detail, movieURL string) error {
	return usecase.Called(ID, title, speaker, detail, movieURL).Error(0)
}

func (usecase *mockPostUseCase) DeletePost(id int) error {
	return usecase.Called(id).Error(0)
}

func getMockPost(id int) *model.Post {
	return &model.Post{
		ID:       id,
		UserID:   id,
		Title:    fmt.Sprintf("title%d", id),
		Speaker:  fmt.Sprintf("speaker%d", id),
		Detail:   fmt.Sprintf("detail%d", id),
		MovieURL: fmt.Sprintf("http://www.example.com/%d", id),
	}
}

func getMockGetPostResult(id int) *model.GetPostResult {
	return &model.GetPostResult{
		Post:     *getMockPost(id),
		UserName: fmt.Sprintf("testuser%d", id),
	}
}

// 登録テスト
func TestCreatePost_success(t *testing.T) {
	// 1. Setup
	post := getMockPost(1)
	jsonBytes, err := json.Marshal(post)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/posts", strings.NewReader(string(jsonBytes)), rec)

	usecase := mockPostUseCase{}
	usecase.On("CreatePost", post.UserID, post.Title, post.Speaker, post.Detail, post.MovieURL).Return(nil)
	handler := NewPostHandler(&usecase)

	// 2. Exercise
	err = handler.CreatePost(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

func TestCreatePost_error_validationError(t *testing.T) {
	cases := []struct {
		label    string
		userID   int
		title    string
		speaker  string
		detail   string
		movieURL string
	}{
		{"UserID下限", 0, "title1", "speaker1", "detail1", "http://example.com/1"},
	}

	for _, test := range cases {
		// 1. Setup
		post := getMockPost(1)
		post.UserID = test.userID
		post.Title = test.title
		post.Speaker = test.speaker
		post.Detail = test.detail
		post.MovieURL = test.movieURL
		jsonBytes, err := json.Marshal(post)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		c := createContext(echo.POST, "/posts", strings.NewReader(string(jsonBytes)), rec)

		usecase := mockPostUseCase{}
		handler := NewPostHandler(&usecase)

		// 2. Exercise
		err = handler.CreatePost(c)

		// 3. Verify
		assert.NoError(t, err, test.label)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)

		// 4. Teardown
	}
}

func TestCreatePost_error_usecaseError(t *testing.T) {
	// 1. Setup
	post := getMockPost(1)
	jsonBytes, err := json.Marshal(post)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/posts", strings.NewReader(string(jsonBytes)), rec)

	usecase := mockPostUseCase{}
	usecase.On("CreatePost", post.UserID, post.Title, post.Speaker, post.Detail, post.MovieURL).Return(errors.New("error"))
	handler := NewPostHandler(&usecase)

	// 2. Exercise
	err = handler.CreatePost(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// 一覧取得テスト
func TestGetPosts_success(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	q := make(url.Values)
	q.Set("limit", "1")
	q.Set("page", "1")
	q.Set("keyword", "")
	c := createContext(echo.GET, "/posts?"+q.Encode(), nil, rec)

	usecase := mockPostUseCase{}
	expected := []*model.GetPostResult{getMockGetPostResult(1), getMockGetPostResult(2)}
	usecase.On("GetPosts", 1, 1, "").Return(2, expected, nil)
	handler := NewPostHandler(&usecase)

	// 2. Exercise
	err := handler.GetPosts(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// TODO 内容検証

	// 4. Teardown
}

func TestGetPosts_error_validationError(t *testing.T) {
	cases := []struct {
		label string
		limit string
		page  string
	}{
		{"limit空", "", "1"},
		{"limit形式", "a", "1"},
		{"limit下限", "0", "1"},
		{"page空", "1", ""},
		{"page形式", "1", "a"},
		{"page下限", "1", "0"},
	}

	for _, test := range cases {
		// 1. Setup
		rec := httptest.NewRecorder()
		q := make(url.Values)
		q.Set("limit", test.limit)
		q.Set("page", test.page)
		q.Set("keyword", "")
		c := createContext(echo.GET, "/posts?"+q.Encode(), nil, rec)

		usecase := mockPostUseCase{}
		handler := NewPostHandler(&usecase)

		// 2. Exercise
		err := handler.GetPosts(c)

		// 3. Verify
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		// 4. Teardown
	}
}

// TODO ここから
// func TestLogin_error_usecaseError(t *testing.T) {
// 	// 1. Setup
// 	email := "testpost@example.com"
// 	password := "testpost"
// 	reader := strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password))
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.POST, "/posts", reader, rec)

// 	usecase := mockPostUseCase{}
// 	usecase.On("Login", email, password).Return(0, "", errors.New("error"))
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err := handler.Login(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)

// 	// 4. Teardown
// }

// // ユーザー詳細テスト
// func TestGetPost_success(t *testing.T) {
// 	// 1. Setup
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.GET, "/posts", nil, rec)
// 	c.SetPath("/posts/:id")
// 	c.SetParamNames("id")
// 	id := 1
// 	c.SetParamValues(fmt.Sprint(id))

// 	expectedPost := getMockPost(id)

// 	usecase := mockPostUseCase{}
// 	usecase.On("GetPost", id).Return(expectedPost, nil)
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err := handler.GetPost(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	post := &model.Post{}
// 	json.Unmarshal(rec.Body.Bytes(), post)
// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	assert.Equal(t, expectedPost, post)

// 	// 4. Teardown
// }

// func TestGetPost_error_validationError(t *testing.T) {
// 	cases := []struct {
// 		label   string
// 		id      interface{}
// 		message string
// 	}{
// 		{"必須", "", "\"ID：数値で入力してください。\"\n"},
// 		{"型", "a", "\"ID：数値で入力してください。\"\n"},
// 		{"下限", 0, "\"ID：1以上の値を入力してください。\"\n"},
// 	}

// 	for _, test := range cases {
// 		// 1. Setup
// 		rec := httptest.NewRecorder()
// 		c := createContext(echo.GET, "/posts", nil, rec)
// 		c.SetPath("/posts/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues(fmt.Sprint(test.id))

// 		usecase := mockPostUseCase{}
// 		handler := NewPostHandler(&usecase)

// 		// 2. Exercise
// 		err := handler.GetPost(c)

// 		// 3. Verify
// 		assert.NoError(t, err, test.label)
// 		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
// 		assert.Equal(t, test.message, rec.Body.String(), test.label)

// 		// 4. Teardown
// 	}
// }

// func TestGetPost_error_usecaseError(t *testing.T) {
// 	// 1. Setup
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.GET, "/posts", nil, rec)
// 	c.SetPath("/posts/:id")
// 	c.SetParamNames("id")
// 	id := 1
// 	c.SetParamValues(fmt.Sprint(id))

// 	usecase := mockPostUseCase{}
// 	usecase.On("GetPost", id).Return(nil, errors.New("error"))
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err := handler.GetPost(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)

// 	// 4. Teardown
// }

// func TestUpdatePost_success(t *testing.T) {
// 	// 1. Setup
// 	post := getMockPost(1)
// 	jsonBytes, err := json.Marshal(post)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.PUT, "/posts", strings.NewReader(string(jsonBytes)), rec)
// 	c.SetPath("/posts/:id")
// 	c.SetParamNames("id")
// 	c.SetParamValues(fmt.Sprint(1))

// 	usecase := mockPostUseCase{}
// 	usecase.On("UpdatePost", post.ID, post.Name, post.Email, post.Password, post.ImageURL).Return(nil)
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err = handler.UpdatePost(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, rec.Code)

// 	// 4. Teardown
// }

// func TestUpdatePost_error_validationError(t *testing.T) {
// 	cases := []struct {
// 		label   string
// 		id      interface{}
// 		name    string
// 		email   string
// 		message string
// 	}{
// 		{"ID空", "", "testpost", "testpost@example.com", "\"ID：数値で入力してください。\"\n"},
// 		{"ID形式", "a", "testpost", "testpost@example.com", "\"ID：数値で入力してください。\"\n"},
// 		{"ID下限", 0, "testpost", "testpost@example.com", "\"ID：必須です。\"\n"},
// 		{"Name空", 1, "", "testpost@example.com", "\"Name：必須です。\"\n"},
// 		{"Email空", 1, "testpost", "", "\"Email：必須です。\"\n"},
// 		{"Email形式", 1, "testpost", "testpostexample.com", "\"Email：正しい形式で入力してください。\"\n"},
// 	}

// 	for _, test := range cases {
// 		// 1. Setup
// 		rec := httptest.NewRecorder()
// 		c := createContext(echo.PUT, "/posts", strings.NewReader(fmt.Sprintf(`{
// 			"name": "%s",
// 			"email": "%s",
// 			"password": "testpost"
// 		}`, test.name, test.email)), rec)
// 		c.SetPath("/posts/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues(fmt.Sprint(test.id))

// 		usecase := mockPostUseCase{}
// 		handler := NewPostHandler(&usecase)

// 		// 2. Exercise
// 		err := handler.UpdatePost(c)

// 		// 3. Verify
// 		assert.NoError(t, err)
// 		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
// 		assert.Equal(t, test.message, rec.Body.String(), test.label)

// 		// 4. Teardown
// 	}
// }

// func TestUpdatePost_error_usecaseError(t *testing.T) {
// 	// 1. Setup
// 	post := getMockPost(1)
// 	jsonBytes, err := json.Marshal(post)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.PUT, "/posts", strings.NewReader(string(jsonBytes)), rec)
// 	c.SetPath("/posts/:id")
// 	c.SetParamNames("id")
// 	id := 1
// 	c.SetParamValues(fmt.Sprint(id))

// 	usecase := mockPostUseCase{}
// 	usecase.On("UpdatePost", post.ID, post.Name, post.Email, post.Password, post.ImageURL).Return(errors.New("error"))
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err = handler.UpdatePost(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)

// 	// 4. Teardown
// }

// // ユーザー削除テスト
// func TestDeletePost_success(t *testing.T) {
// 	// 1. Setup
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.DELETE, "/posts", nil, rec)
// 	c.SetPath("/posts/:id")
// 	c.SetParamNames("id")
// 	id := 1
// 	c.SetParamValues(fmt.Sprint(id))

// 	usecase := mockPostUseCase{}
// 	usecase.On("DeletePost", id).Return(nil)
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err := handler.DeletePost(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, rec.Code)

// 	// 4. Teardown
// }

// func TestDeletePost_error_validationError(t *testing.T) {
// 	cases := []struct {
// 		label   string
// 		id      interface{}
// 		message string
// 	}{
// 		{"ID空", "", "\"ID：数値で入力してください。\"\n"},
// 		{"ID形式", "a", "\"ID：数値で入力してください。\"\n"},
// 		{"ID下限", 0, "\"ID：1以上の値を入力してください。\"\n"},
// 	}

// 	for _, test := range cases {
// 		// 1. Setup
// 		rec := httptest.NewRecorder()
// 		c := createContext(echo.DELETE, "/posts", nil, rec)
// 		c.SetPath("/posts/:id")
// 		c.SetParamNames("id")
// 		c.SetParamValues(fmt.Sprint(test.id))

// 		usecase := mockPostUseCase{}
// 		handler := NewPostHandler(&usecase)

// 		// 2. Exercise
// 		err := handler.DeletePost(c)

// 		// 3. Verify
// 		assert.NoError(t, err, test.label)
// 		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
// 		assert.Equal(t, test.message, rec.Body.String(), test.label)

// 		// 4. Teardown
// 	}
// }

// func TestDeletePost_error_usecaseError(t *testing.T) {
// 	// 1. Setup
// 	rec := httptest.NewRecorder()
// 	c := createContext(echo.DELETE, "/posts", nil, rec)
// 	c.SetPath("/posts/:id")
// 	c.SetParamNames("id")
// 	id := 1
// 	c.SetParamValues(fmt.Sprint(id))

// 	usecase := mockPostUseCase{}
// 	usecase.On("DeletePost", id).Return(errors.New("error"))
// 	handler := NewPostHandler(&usecase)

// 	// 2. Exercise
// 	err := handler.DeletePost(c)

// 	// 3. Verify
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)

// 	// 4. Teardown
// }
