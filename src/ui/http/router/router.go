package router

import (
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/handler"
	"github.com/labstack/echo"
)

// SetRoutes Routerの設定を行います.
func SetRoutes(e *echo.Echo, handler handler.AppHandler) {
	e.POST("/users", handler.CreateUser)
	e.POST("/login", handler.Login)
	e.GET("/users", handler.GetUsers)
	e.GET("/users/:id", handler.GetUser)
	e.PUT("/users/:id", handler.UpdateUser)
	e.DELETE("/users/:id", handler.DeleteUser)

	// TODO 以下削除
	// e.POST("/signup", handler.Signup)
	// e.POST("/login", handler.Login)
	// e.GET("/api/v1/posts", model.GetAllPosts)
	// e.GET("/api/v1/posts/:id", model.GetPost)

	// api := e.Group("/api/v1")
	// api.Use(middleware.JWTWithConfig(handler.Config))
	// api.GET("/users", model.GetAllUsers)
}
