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

func (repository *mockPostRepository) Fetch(limit, page int, keyword string) (int, []*model.GetPostResult, error) {
	args := repository.Called(limit, page, keyword)
	posts, ok := args.Get(1).([]*model.GetPostResult)
	if ok {
		return args.Int(0), posts, args.Error(2)
	} else {
		return args.Int(0), nil, args.Error(2)
	}
}

func (repository *mockPostRepository) FetchByID(id int) (*model.GetPostResult, error) {
	args := repository.Called(id)
	post, ok := args.Get(0).(*model.GetPostResult)
	if ok {
		return post, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
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
		MovieURL: fmt.Sprintf("https://www.example.com/watch?v=%d", id),
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

func getMockGetPostResult(id int) *model.GetPostResult {
	post := getMockPostForRead(id)
	return &model.GetPostResult{
		Post:     *post,
		UserName: fmt.Sprintf("username%d", id),
	}
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
	expectedTotalCount := 2
	expectedPosts := []*model.GetPostResult{getMockGetPostResult(1), getMockGetPostResult(2)}
	repository.On("Fetch", limit, page, keyword).Return(expectedTotalCount, expectedPosts, nil)

	// 2. Exercise
	totalCount, posts, err := usecase.GetPosts(limit, page, keyword)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedTotalCount, totalCount)
	assert.Equal(t, len(expectedPosts), len(posts))
	assert.Equal(t, expectedPosts[0], posts[0])
	assert.Equal(t, expectedPosts[1], posts[1])

	// 4. Teardown
}

func TestGetPosts_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	limit := 3
	page := 1
	keyword := ""
	repository.On("Fetch", limit, page, keyword).Return(0, nil, errors.New("error"))

	// 2. Execise
	totalCount, posts, err := usecase.GetPosts(limit, page, keyword)

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, totalCount)
	assert.Empty(t, posts)

	// 4. Teardown
}

// 投稿詳細テスト
func TestGetPost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	expected := getMockGetPostResult(id)
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
	assert.Equal(t, expected.UserName, post.UserName)

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
