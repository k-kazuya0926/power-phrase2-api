package interactor

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"github.com/k-kazuya0926/power-phrase2-api/domain/service"
	"github.com/k-kazuya0926/power-phrase2-api/infrastructure/persistence/datastore"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/handler"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
)

// Interactor interfase Intractorは簡易DIコンテナとしての役割を持つ.
type Interactor interface {
	// TODO 削除
	// NewUserRepository() repository.UserRepository
	// NewUserService() service.UserService
	// NewUserUseCase() usecase.UserUseCase
	// NewUserHandler() handler.UserHandler
	NewAppHandler() handler.AppHandler
}

type interactor struct {
}

// NewInteractor intractorを取得します.
func NewInteractor() Interactor {
	return &interactor{}
}

func (interactor *interactor) NewAppHandler() handler.AppHandler {
	return handler.NewAppHandler(interactor.NewUserHandler(), interactor.NewPostHandler())
}

// User
func (interactor *interactor) NewUserRepository() repository.UserRepository {
	return datastore.NewUserRepository()
}

func (interactor *interactor) NewUserService() service.UserService {
	return service.NewUserService(interactor.NewUserRepository())
}

func (interactor *interactor) NewUserUseCase() usecase.UserUseCase {
	return usecase.NewUserUseCase(interactor.NewUserRepository())
}

func (interactor *interactor) NewUserHandler() handler.UserHandler {
	return handler.NewUserHandler(interactor.NewUserUseCase())
}

// Post
func (interactor *interactor) NewPostRepository() repository.PostRepository {
	return datastore.NewPostRepository()
}

func (interactor *interactor) NewPostUseCase() usecase.PostUseCase {
	return usecase.NewPostUseCase(interactor.NewPostRepository())
}

func (interactor *interactor) NewPostHandler() handler.PostHandler {
	return handler.NewPostHandler(interactor.NewPostUseCase())
}
