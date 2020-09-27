package main

import (
	"app/handler"
	"app/model"

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

	prefix := "/api/v1/"
	e.GET(prefix+"users", model.GetAllUsers)
	e.GET(prefix+"user/:id", model.GetUser)
	e.PUT(prefix+"user/:id", model.UpdateUser)
	e.DELETE(prefix+"user/:id", model.DeleteUser)

	// e.GET(prefix+"posts", model.GetAllPosts)
	// e.GET(prefix+"posts/:id", model.GetPost)
	// e.POST(prefix+"posts", model.CreatePost)
	// e.PUT(prefix+"posts/:id", model.UpdatePost)
	// e.DELETE(prefix+"posts/:id", model.DeletePost)
}
