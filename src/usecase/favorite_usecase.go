// Package usecase Application Service層。
package usecase

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// FavoriteUseCase インターフェース
type FavoriteUseCase interface {
	CreateFavorite(postID, userID int) (err error)
	GetFavorites(userID, limit, offset int) (totalCount int, posts []*model.GetPostResult, err error)
	DeleteFavorite(id int) error
}

// favoriteUseCase 構造体
type favoriteUseCase struct {
	repository.PostRepository
}

// NewFavoriteUseCase FavoriteUseCaseを生成。
func NewFavoriteUseCase(repository repository.PostRepository) FavoriteUseCase {
	return &favoriteUseCase{repository}
}

// CreateFavorite お気に入り登録
func (usecase *favoriteUseCase) CreateFavorite(userID, postID int) (err error) {
	favorite := model.Favorite{
		UserID: userID,
		PostID: postID,
	}
	err = usecase.PostRepository.CreateFavorite(&favorite)

	return err
}

// GetFavorites お気に入り一覧取得
func (usecase *favoriteUseCase) GetFavorites(userID, limit, page int) (totalCount int, posts []*model.GetPostResult, err error) {
	totalCount, posts, err = usecase.PostRepository.FetchFavorites(userID, limit, page)
	if err != nil {
		return 0, nil, err
	}

	return totalCount, posts, nil
}

// DeleteFavorite お気に入り削除
func (usecase *favoriteUseCase) DeleteFavorite(id int) error {
	if err := usecase.PostRepository.DeleteFavorite(id); err != nil {
		return err
	}
	return nil
}
