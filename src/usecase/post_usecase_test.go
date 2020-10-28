package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock
type mockPostRepository struct {
	mock.Mock
}

func (repository *mockPostRepository) Create(post *model.Post) error {
	return repository.Called(post).Error(0)
}

func (repository *mockPostRepository) Fetch(limit, page int, keyword string) ([]*model.Post, error) {
	args := repository.Called(limit, page, keyword)
	posts := args.Get(0)
	if posts == nil {
		return nil, args.Error(1)
	}
	return posts.([]*model.Post), args.Error(1)
}

func (repository *mockPostRepository) FetchByID(id int) (*model.Post, error) {
	args := repository.Called(id)
	post := args.Get(0)
	if post == nil {
		return nil, args.Error(1)
	}
	return post.(*model.Post), args.Error(1)
}

func (repository *mockPostRepository) Update(post *model.Post) error {
	return repository.Called(post).Error(0)
}

func (repository *mockPostRepository) Delete(id int) error {
	return repository.Called(id).Error(0)
}

// 入力用投稿
func getMockPostForInput(id int) *model.Post {
	post := &model.Post{
		ID:       id,
		UserID:   id,
		Title:    fmt.Sprintf("title%d", id),
		Speaker:  fmt.Sprintf("speaker%d", id),
		Detail:   fmt.Sprintf("detail%d", id),
		MovieURL: fmt.Sprintf("http://www.example.com/%d", id),
	}
	return post
}

// DBから取得された投稿
func getMockPostForRead(id int) *model.Post {
	post := getMockPostForInput(id)
	post.CreatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	post.UpdatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	return post
}

// 投稿登録テスト
func TestCreatePost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := getMockPostForInput(id)
	repository.On("Create", mock.AnythingOfType("*model.Post")).Return(nil)

	// 2. Exercise
	err := usecase.CreatePost(post.UserID, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestCreatePost_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := getMockPostForInput(id)
	repository.On("Create", mock.AnythingOfType("*model.Post")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.CreatePost(post.UserID, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// 投稿一覧テスト
func TestGetPosts_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	limit := 3
	page := 1
	keyword := ""
	expected := []*model.Post{getMockPostForRead(1), getMockPostForRead(2)}
	repository.On("Fetch", limit, page, keyword).Return(expected, nil)

	// 2. Exercise
	posts, err := usecase.GetPosts(limit, page, keyword)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(posts))
	assert.Equal(t, expected[0], posts[0])
	assert.Equal(t, expected[1], posts[1])

	// 4. Teardown
}

func TestGetPosts_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	limit := 3
	page := 1
	keyword := ""
	repository.On("Fetch", limit, page, keyword).Return(nil, errors.New("error"))

	// 2. Execise
	posts, err := usecase.GetPosts(limit, page, keyword)

	// 3. Verify
	assert.Error(t, err)
	assert.Empty(t, posts)

	// 4. Teardown
}

// 投稿詳細テスト
func TestGetPost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	expected := getMockPostForRead(id)
	repository.On("FetchByID", id).Return(expected, nil)

	// 2. Exercise
	post, err := usecase.GetPost(id)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, post.ID)
	assert.Equal(t, expected.CreatedAt, post.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, post.UpdatedAt)
	assert.Equal(t, expected.DeletedAt, post.DeletedAt)
	assert.Equal(t, expected.UserID, post.UserID)
	assert.Equal(t, expected.Title, post.Title)
	assert.Equal(t, expected.Speaker, post.Speaker)
	assert.Equal(t, expected.Detail, post.Detail)
	assert.Equal(t, expected.MovieURL, post.MovieURL)

	// 4. Teardown
}

func TestGetPost_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	repository.On("FetchByID", id).Return(nil, errors.New("error"))

	// 2. Execise
	post, err := usecase.GetPost(id)

	// 3. Verify
	assert.Error(t, err)
	assert.Empty(t, post)

	// 4. Teardown
}

// 投稿更新テスト
func TestUpdatePost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := getMockPostForInput(id)
	repository.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)

	// 2. Exercise
	err := usecase.UpdatePost(id, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestUpdatePost_error(t *testing.T) {
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := getMockPostForInput(id)
	repository.On("Update", mock.AnythingOfType("*model.Post")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.UpdatePost(id, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// 投稿削除テスト
func TestDeletePost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(nil)

	// 2. Exercise
	err := usecase.DeletePost(id)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestDeletePost_error(t *testing.T) {
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.DeletePost(id)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}