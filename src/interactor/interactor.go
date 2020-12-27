// Package interactor 簡易DIコンテナ
package interactor

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"github.com/k-kazuya0926/power-phrase2-api/infrastructure/persistence/datastore"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/handler"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
)

// Interactor インターフェース。AppHandlerのインターフェースを保持。
type Interactor interface {
	NewAppHandler() handler.AppHandler
}

// interactor 構造体
type interactor struct {
}

// NewInteractor intractorを生成。
func NewInteractor() Interactor {
	return &interactor{}
}

// NewAppHandler AppHandlerを生成。
func (interactor *interactor) NewAppHandler() handler.AppHandler {
	return handler.NewAppHandler(interactor.NewUserHandler(), interactor.NewPostHandler(), interactor.NewCommentHandler(), interactor.NewFavoriteHandler())
}

// ユーザー関連
// NewUserRepository UserRepositoryを生成。
func (interactor *interactor) NewUserRepository() repository.UserRepository {
	return datastore.NewUserRepository()
}

// NewUserUseCase UserUseCaseを生成。
func (interactor *interactor) NewUserUseCase() usecase.UserUseCase {
	return usecase.NewUserUseCase(interactor.NewUserRepository())
}

// NewUserHandler UserHandlerを生成。
func (interactor *interactor) NewUserHandler() handler.UserHandler {
	return handler.NewUserHandler(interactor.NewUserUseCase())
}

// 投稿関連
// NewPostRepository PostRepositoryを生成。
func (interactor *interactor) NewPostRepository() repository.PostRepository {
	return datastore.NewPostRepository()
}

// NewPostUseCase PostUseCaseを生成。
func (interactor *interactor) NewPostUseCase() usecase.PostUseCase {
	return usecase.NewPostUseCase(interactor.NewPostRepository())
}

// NewPostHandler PostHandlerを生成。
func (interactor *interactor) NewPostHandler() handler.PostHandler {
	return handler.NewPostHandler(interactor.NewPostUseCase())
}

// コメント関連
// NewCommentUseCase CommentUseCaseを生成。
func (interactor *interactor) NewCommentUseCase() usecase.CommentUseCase {
	return usecase.NewCommentUseCase(interactor.NewPostRepository())
}

// NewCommentHandler CommentHandlerを生成。
func (interactor *interactor) NewCommentHandler() handler.CommentHandler {
	return handler.NewCommentHandler(interactor.NewCommentUseCase())
}

// お気に入り関連
// NewFavoriteUseCase FavoriteUseCaseを生成。
func (interactor *interactor) NewFavoriteUseCase() usecase.FavoriteUseCase {
	return usecase.NewFavoriteUseCase(interactor.NewPostRepository())
}

// NewFavoriteHandler FavoriteHandlerを生成。
func (interactor *interactor) NewFavoriteHandler() handler.FavoriteHandler {
	return handler.NewFavoriteHandler(interactor.NewFavoriteUseCase())
}
