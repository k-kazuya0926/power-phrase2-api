// Package handler UI層
package handler

// AppHandler 全てのHandlerのinterfaceを満たす。
type AppHandler interface {
	UserHandler
	PostHandler
	// embed all handler interfaces
}

// appHandler 構造体
type appHandler struct {
	UserHandler
	PostHandler
	// embed all handler interfaces
}

// NewAppHandler AppHandlerを生成
func NewAppHandler(userHandler UserHandler, postHandler PostHandler) AppHandler {
	return &appHandler{userHandler, postHandler}
}
