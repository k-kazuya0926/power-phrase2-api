package router

import (
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/handler"
	"github.com/labstack/echo"
)

// SetRoutes Routerの設定を行います.
func SetRoutes(e *echo.Echo, handler handler.AppHandler) {
	e.POST("/users", handler.CreateUser)
	e.POST("/login", handler.Login)
	e.GET("/users/:id", handler.GetUser)
	e.PUT("/users/:id", handler.UpdateUser)
	e.DELETE("/users/:id", handler.DeleteUser)

	e.POST("/posts", handler.CreatePost)
	// e.GET("/api/v1/posts", handler.GetAllPosts)
	e.GET("/posts/:id", handler.GetPost)
	e.PUT("/posts/:id", handler.UpdatePost)
	e.DELETE("/posts/:id", handler.DeletePost)
}
