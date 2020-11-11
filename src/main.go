package main

import (
	"fmt"
	"os"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/interactor"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/router"
	"github.com/k-kazuya0926/power-phrase2-api/validator"
	"github.com/labstack/echo"
)

func main() {
	conf.NewConfig()

	e := echo.New()
	interactor := interactor.NewInteractor() // TODO DI用のライブラリはないのかな？
	handler := interactor.NewAppHandler()

	router.SetRoutes(e, handler)

	e.Validator = validator.NewValidator()

	if err := e.Start(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))); err != nil {
		e.Logger.Fatal(fmt.Sprintf("Failed to start: %v", err))
	}
}
