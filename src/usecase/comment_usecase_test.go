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

// 入力用コメント
func makeCommentForInput(id, postID, userID int) *model.Comment {
	comment := &model.Comment{
		ID:     id,
		PostID: postID,
		UserID: userID,
		Body:   fmt.Sprintf("body%d", id),
	}
	return comment
}

// DBから取得されたコメント
func makeCommentForRead(id, postID, userID int) *model.Comment {
	comment := makeCommentForInput(id, postID, userID)
	comment.CreatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	comment.UpdatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	return comment
}

func makeGetCommentResult(id, postID, userID int) *model.GetCommentResult {
	comment := makeCommentForRead(id, postID, userID)
	return &model.GetCommentResult{
		Comment:           *comment,
		UserName:          fmt.Sprintf("username%d", id),
		UserImageFilePath: fmt.Sprintf("images/%d.png", id),
	}
}

// コメント登録成功
func TestCreateComment_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	postID := 1
	userID := 1
	comment := makeCommentForInput(id, postID, userID)
	repository.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	// 2. Exercise
	err := usecase.CreateComment(comment.PostID, comment.UserID, comment.Body)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

// コメント登録エラー
func TestCreateComment_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	postID := 1
	userID := 1
	comment := makeCommentForInput(id, postID, userID)
	repository.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.CreateComment(comment.PostID, comment.UserID, comment.Body)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// コメント一覧取得成功
func TestGetComments_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewCommentUseCase(&repository)
	limit := 3
	page := 1
	postID := 1
	expectedTotalCount := 2
	expectedComments := []*model.GetCommentResult{makeGetCommentResult(1, postID, 1), makeGetCommentResult(2, postID, 2)}
	repository.On("FetchComments", postID, limit, page).Return(expectedTotalCount, expectedComments, nil)

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

// コメント一覧取得エラー
func TestGetComments_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewCommentUseCase(&repository)
	limit := 3
	page := 1
	postID := 1
	repository.On("FetchComments", postID, limit, page).Return(0, nil, errors.New("error"))

	// 2. Execise
	totalCount, comments, err := usecase.GetComments(postID, limit, page)

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, totalCount)
	assert.Empty(t, comments)

	// 4. Teardown
}

// コメント削除成功
func TestDeleteComment_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	repository.On("DeleteComment", id).Return(nil)

	// 2. Exercise
	err := usecase.DeleteComment(id)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

// コメント削除エラー
func TestDeleteComment_error(t *testing.T) {
	repository := mockPostRepository{}
	usecase := NewCommentUseCase(&repository)
	id := 1
	repository.On("DeleteComment", id).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.DeleteComment(id)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}
