package main

import (
	"fmt"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/interactor"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/router"
	"github.com/k-kazuya0926/power-phrase2-api/validator"
	"github.com/labstack/echo"
)

func main() {
	conf.NewConfig(false)

	e := echo.New()
	interactor := interactor.NewInteractor() // TODO DI用のライブラリはないのかな？
	handler := interactor.NewAppHandler()

	router.SetRoutes(e, handler)

	e.Validator = validator.NewValidator()

	if err := e.Start(fmt.Sprintf(":%d", conf.Current.Server.Port)); err != nil {
		e.Logger.Fatal(fmt.Sprintf("Failed to start: %v", err))
	}
}
