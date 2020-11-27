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
type mockCommentRepository struct {
	mock.Mock
}

func (repository *mockCommentRepository) Create(comment *model.Comment) error {
	return repository.Called(comment).Error(0)
}

func (repository *mockCommentRepository) Fetch(postID, limit, page int) (int, []*model.GetCommentResult, error) {
	args := repository.Called(postID, limit, page)
	comments, ok := args.Get(1).([]*model.GetCommentResult)
	if ok {
		return args.Int(0), comments, args.Error(2)
	} else {
		return args.Int(0), nil, args.Error(2)
	}
}

func (repository *mockCommentRepository) Delete(id int) error {
	return repository.Called(id).Error(0)
}

// 入力用コメント
func getMockCommentForInput(id, postID, userID int) *model.Comment {
	comment := &model.Comment{
		ID:     id,
		PostID: postID,
		UserID: userID,
		Body:   fmt.Sprintf("body%d", id),
	}
	return comment
}

// DBから取得されたコメント
func getMockCommentForRead(id, postID, userID int) *model.Comment {
	comment := getMockCommentForInput(id, postID, userID)
	comment.CreatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	comment.UpdatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	return comment
}

func getMockGetCommentResult(id, postID, userID int) *model.GetCommentResult {
	comment := getMockCommentForRead(id, postID, userID)
	return &model.GetCommentResult{
		Comment:           *comment,
		UserName:          fmt.Sprintf("username%d", id),
		UserImageFilePath: fmt.Sprintf("images/%d.png", id),
	}
}

// コメント登録テスト
func TestCreateComment_success(t *testing.T) {
	// 1. Setup
	repository := mockCommentRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	postID := 1
	userID := 1
	comment := getMockCommentForInput(id, postID, userID)
	repository.On("Create", mock.AnythingOfType("*model.Comment")).Return(nil)

	// 2. Exercise
	err := usecase.CreateComment(comment.PostID, comment.UserID, comment.Body)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestCreateComment_error(t *testing.T) {
	// 1. Setup
	repository := mockCommentRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	postID := 1
	userID := 1
	comment := getMockCommentForInput(id, postID, userID)
	repository.On("Create", mock.AnythingOfType("*model.Comment")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.CreateComment(comment.PostID, comment.UserID, comment.Body)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// コメント一覧テスト
func TestGetComments_success(t *testing.T) {
	// 1. Setup
	repository := mockCommentRepository{}
	usecase := NewCommentUseCase(&repository)
	limit := 3
	page := 1
	postID := 1
	expectedTotalCount := 2
	expectedComments := []*model.GetCommentResult{getMockGetCommentResult(1, postID, 1), getMockGetCommentResult(2, postID, 2)}
	repository.On("Fetch", postID, limit, page).Return(expectedTotalCount, expectedComments, nil)

	// 2. Exercise
	totalCount, comments, err := usecase.GetComments(postID, limit, page)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedTotalCount, totalCount)
	assert.Equal(t, len(expectedComments), len(comments))
	assert.Equal(t, expectedComments[0], comments[0])
	assert.Equal(t, expectedComments[1], comments[1])

	// 4. Teardown
}

func TestGetComments_error(t *testing.T) {
	// 1. Setup
	repository := mockCommentRepository{}
	usecase := NewCommentUseCase(&repository)
	limit := 3
	page := 1
	postID := 1
	repository.On("Fetch", postID, limit, page).Return(0, nil, errors.New("error"))

	// 2. Execise
	totalCount, comments, err := usecase.GetComments(postID, limit, page)

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, totalCount)
	assert.Empty(t, comments)

	// 4. Teardown
}

// コメント削除テスト
func TestDeleteComment_success(t *testing.T) {
	// 1. Setup
	repository := mockCommentRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(nil)

	// 2. Exercise
	err := usecase.DeleteComment(id)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestDeleteComment_error(t *testing.T) {
	repository := mockCommentRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.DeleteComment(id)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}
