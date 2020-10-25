package handler

// AppHandler interfase
// AppHandlerは全てのHandlerのinterfaceを満たす.※routerの実装が依存する.
type AppHandler interface {
	UserHandler
	PostHandler
	// embed all handler interfaces
}

type appHandler struct {
	UserHandler
	PostHandler
	// embed all handler interfaces
}

func NewAppHandler(userHandler UserHandler, postHandler PostHandler) AppHandler {
	return &appHandler{userHandler, postHandler}
}
