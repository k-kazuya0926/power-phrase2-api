// Package usecase Application Service層。
package usecase

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// CommentUseCase インターフェース
type CommentUseCase interface {
	CreateComment(postID, userID int, body string) (err error)
	GetComments(postID, limit, offset int) (totalCount int, comments []*model.GetCommentResult, err error)
	DeleteComment(id int) error
}

// commentUseCase 構造体
type commentUseCase struct {
	repository.CommentRepository
}

// NewCommentUseCase CommentUseCaseを生成。
func NewCommentUseCase(repository repository.CommentRepository) CommentUseCase {
	return &commentUseCase{repository}
}

// CreateComment 登録
func (usecase *commentUseCase) CreateComment(postID, userID int, body string) (err error) {
	comment := model.Comment{
		PostID: postID,
		UserID: userID,
		Body:   body,
	}
	err = usecase.CommentRepository.Create(&comment)

	return err
}

// GetComments 一覧取得
func (usecase *commentUseCase) GetComments(postID, limit, page int) (totalCount int, comments []*model.GetCommentResult, err error) {
	totalCount, comments, err = usecase.CommentRepository.Fetch(postID, limit, page)
	if err != nil {
		return 0, nil, err
	}

	return totalCount, comments, nil
}

// DeleteComment 削除
func (usecase *commentUseCase) DeleteComment(id int) error {
	if err := usecase.CommentRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
