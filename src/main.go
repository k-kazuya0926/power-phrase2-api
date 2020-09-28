package main

import (
	"github.com/k-kazuya0926/power-phrase2-api/handler"
	"github.com/k-kazuya0926/power-phrase2-api/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	initRouting(e)
	e.Logger.Fatal(e.Start(":9000"))
}

func initRouting(e *echo.Echo) {
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	api := e.Group("/api/v1")
	api.Use(middleware.JWTWithConfig(handler.Config))
	api.GET("/users", model.GetAllUsers)
	api.GET("/user/:id", model.GetUser)
	api.PUT("/user/:id", model.UpdateUser)
	api.DELETE("/user/:id", model.DeleteUser)

	// api.GET("posts", model.GetAllPosts)
	// api.GET("posts/:id", model.GetPost)
	// api.POST("posts", model.CreatePost)
	// api.PUT("posts/:id", model.UpdatePost)
	// api.DELETE("posts/:id", model.DeletePost)
}
