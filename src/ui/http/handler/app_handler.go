// Package handler UI層
package handler

// AppHandler 全てのHandlerのinterfaceを満たす。
type AppHandler interface {
	UserHandler
	PostHandler
	CommentHandler
	// embed all handler interfaces
}

// appHandler 構造体
type appHandler struct {
	UserHandler
	PostHandler
	CommentHandler
	// embed all handler interfaces
}

// NewAppHandler AppHandlerを生成
func NewAppHandler(userHandler UserHandler, postHandler PostHandler, commentHandler CommentHandler) AppHandler {
	return &appHandler{userHandler, postHandler, commentHandler}
}
