package main

import (
	"flag"
	"fmt"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/interactor"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/middleware"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/router"
	"github.com/k-kazuya0926/power-phrase2-api/validator"
	"github.com/labstack/echo"
)

// TODO調査
//Dockerコンテナで実行する時(production.confをもとに起動するとき)は起動時に-serverを指定
var runServer = flag.Bool("server", false, "production is -server option require")

func main() {
	flag.Parse()
	conf.NewConfig(*runServer)

	e := echo.New()
	interactor := interactor.NewInteractor() // TODO DI用のライブラリはないのかな？
	handler := interactor.NewAppHandler()

	router.SetRoutes(e, handler)
	middleware.SetMiddlewares(e)

	e.Validator = validator.NewValidator()

	if err := e.Start(fmt.Sprintf(":%d", conf.Current.Server.Port)); err != nil {
		e.Logger.Fatal(fmt.Sprintf("Failed to start: %v", err))
	}
}
