package interactor

import (
	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"github.com/k-kazuya0926/power-phrase2-api/domain/service"
	"github.com/k-kazuya0926/power-phrase2-api/infrastructure/persistence/datastore"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/handler"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
)

// Interactor interfase Intractorは簡易DIコンテナとしての役割を持つ.
type Interactor interface {
	NewUserRepository() repository.UserRepository
	NewUserService() service.UserService
	NewUserUseCase() usecase.UserUseCase
	NewUserHandler() handler.UserHandler
	NewAppHandler() handler.AppHandler
}

type interactor struct {
	Connection *gorm.DB
}

// NewInteractor intractorを取得します.
func NewInteractor(connection *gorm.DB) Interactor {
	return &interactor{connection}
}

type appHandler struct { // TODO なぜapp_handler.goではなくここに定義されているのか？
	handler.UserHandler
	// embed all handler interfaces
}

func (interactor *interactor) NewAppHandler() handler.AppHandler {
	appHandler := &appHandler{}
	appHandler.UserHandler = interactor.NewUserHandler()
	return appHandler
}

func (interactor *interactor) NewUserRepository() repository.UserRepository {
	return datastore.NewUserRepository(interactor.Connection)
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