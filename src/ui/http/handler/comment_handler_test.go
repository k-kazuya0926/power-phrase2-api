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
type mockCommentUseCase struct {
	mock.Mock
}

func (usecase *mockCommentUseCase) CreateComment(postID, userID int, body string) (err error) {
	return usecase.Called(postID, userID, body).Error(0)
}

func (usecase *mockCommentUseCase) GetComments(postID, limit, offset int) (totalCount int, comments []*model.GetCommentResult, err error) {
	args := usecase.Called(postID, limit, offset)
	comments, ok := args.Get(1).([]*model.GetCommentResult)
	if ok {
		return args.Int(0), comments, args.Error(2)
	}

	return args.Int(0), nil, args.Error(2)
}

func (usecase *mockCommentUseCase) DeleteComment(id int) error {
	return usecase.Called(id).Error(0)
}

func getMockComment(id, postID, userID int) *model.Comment {
	return &model.Comment{
		ID:     id,
		PostID: postID,
		UserID: userID,
		Body:   fmt.Sprintf("body%d", id),
	}
}

func getMockGetCommentResult(id, postID, userID int) *model.GetCommentResult {
	return &model.GetCommentResult{
		Comment:           *getMockComment(id, postID, userID),
		UserName:          fmt.Sprintf("testuser%d", id),
		UserImageFilePath: fmt.Sprintf("images/%d.png", id),
	}
}

// 登録テスト
func TestCreateComment_success(t *testing.T) {
	// 1. Setup
	postID := 1
	comment := getMockComment(1, postID, 1)
	jsonBytes, err := json.Marshal(comment)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/posts/:id/comments", strings.NewReader(string(jsonBytes)), rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(postID))

	usecase := mockCommentUseCase{}
	usecase.On("CreateComment", comment.PostID, comment.UserID, comment.Body).Return(nil)
	handler := NewCommentHandler(&usecase)

	// 2. Exercise
	err = handler.CreateComment(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

func TestCreateComment_error_validationError(t *testing.T) {
	cases := []struct {
		label  string
		postID int
		userID int
		body   string
	}{
		{"PostID下限", 0, 1, "body"},
		{"UserID下限", 1, 0, "body"},
	}

	for _, test := range cases {
		// 1. Setup
		comment := getMockComment(1, test.postID, test.userID)
		comment.Body = test.body
		jsonBytes, err := json.Marshal(comment)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		c := createContext(echo.POST, "/posts/:id/comments", strings.NewReader(string(jsonBytes)), rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.postID))

		usecase := mockCommentUseCase{}
		handler := NewCommentHandler(&usecase)

		// 2. Exercise
		err = handler.CreateComment(c)

		// 3. Verify
		assert.NoError(t, err, test.label)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)

		// 4. Teardown
	}
}

func TestCreateComment_error_usecaseError(t *testing.T) {
	// 1. Setup
	postID := 1
	comment := getMockComment(1, postID, 1)
	jsonBytes, err := json.Marshal(comment)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	c := createContext(echo.POST, "/posts/:id/comments", strings.NewReader(string(jsonBytes)), rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(postID))

	usecase := mockCommentUseCase{}
	usecase.On("CreateComment", comment.PostID, comment.UserID, comment.Body).Return(errors.New("error"))
	handler := NewCommentHandler(&usecase)

	// 2. Exercise
	err = handler.CreateComment(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// 一覧取得テスト
func TestGetComments_success(t *testing.T) {
	// 1. Setup
	postID := 1
	rec := httptest.NewRecorder()
	q := make(url.Values)
	q.Set("limit", "10")
	q.Set("page", "1")
	c := createContext(echo.GET, "/posts/:id/comments?"+q.Encode(), nil, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(postID))

	usecase := mockCommentUseCase{}
	expected := []*model.GetCommentResult{getMockGetCommentResult(1, postID, 1), getMockGetCommentResult(2, postID, 2)}
	usecase.On("GetComments", postID, 10, 1).Return(2, expected, nil)
	handler := NewCommentHandler(&usecase)

	// 2. Exercise
	err := handler.GetComments(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	// TODO 内容検証

	// 4. Teardown
}

func TestGetComments_error_validationError(t *testing.T) {
	cases := []struct {
		label  string
		postID string
		limit  string
		page   string
	}{
		{"postID空", "", "1", "1"},
		{"postID形式", "a", "1", "1"},
		{"postID下限", "0", "a", "1"},
		{"limit空", "1", "", "1"},
		{"limit形式", "1", "a", "1"},
		{"limit下限", "1", "0", "1"},
		{"page空", "1", "1", ""},
		{"page形式", "1", "1", "a"},
		{"page下限", "1", "1", "0"},
	}

	for _, test := range cases {
		// 1. Setup
		rec := httptest.NewRecorder()
		q := make(url.Values)
		q.Set("limit", test.limit)
		q.Set("page", test.page)
		c := createContext(echo.GET, "/posts/:id/comments?"+q.Encode(), nil, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.postID))

		usecase := mockCommentUseCase{}
		handler := NewCommentHandler(&usecase)

		// 2. Exercise
		err := handler.GetComments(c)

		// 3. Verify
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		// 4. Teardown
	}
}

func TestGetComments_error_usecaseError(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	q := make(url.Values)
	q.Set("limit", "1")
	q.Set("page", "1")
	c := createContext(echo.GET, "/posts/:id/comments?"+q.Encode(), nil, rec)
	c.SetParamNames("id")
	postID := 1
	c.SetParamValues(fmt.Sprint(postID))

	usecase := mockCommentUseCase{}
	usecase.On("GetComments", postID, 1, 1).Return(0, nil, errors.New("error"))
	handler := NewCommentHandler(&usecase)

	// 2. Exercise
	err := handler.GetComments(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}

// 削除テスト
func TestDeleteComment_success(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.DELETE, "/comments/:id", nil, rec)
	c.SetParamNames("id")
	postID := 1
	c.SetParamValues(fmt.Sprint(postID))

	usecase := mockCommentUseCase{}
	usecase.On("DeleteComment", postID).Return(nil)
	handler := NewCommentHandler(&usecase)

	// 2. Exercise
	err := handler.DeleteComment(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// 4. Teardown
}

func TestDeleteComment_error_validationError(t *testing.T) {
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
		c := createContext(echo.DELETE, "/comments/:id", nil, rec)
		// c.SetPath("/comments/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(test.id))

		usecase := mockCommentUseCase{}
		handler := NewCommentHandler(&usecase)

		// 2. Exercise
		err := handler.DeleteComment(c)

		// 3. Verify
		assert.NoError(t, err, test.label)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, test.label)
		assert.Equal(t, test.message, rec.Body.String(), test.label)

		// 4. Teardown
	}
}

func TestDeleteComment_error_usecaseError(t *testing.T) {
	// 1. Setup
	rec := httptest.NewRecorder()
	c := createContext(echo.DELETE, "/comments/:id", nil, rec)
	// c.SetPath("/comments/:id")
	c.SetParamNames("id")
	postID := 1
	c.SetParamValues(fmt.Sprint(postID))

	usecase := mockCommentUseCase{}
	usecase.On("DeleteComment", postID).Return(errors.New("error"))
	handler := NewCommentHandler(&usecase)

	// 2. Exercise
	err := handler.DeleteComment(c)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// 4. Teardown
}
