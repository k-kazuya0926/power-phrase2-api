package main

import (
	"flag"
	"fmt"

	_ "github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/interactor"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/middleware"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/router"
	"github.com/labstack/echo"
)

//Dockerコンテナで実行する時(production.confをもとに起動するとき)は起動時に-serverを指定
var runServer = flag.Bool("server", false, "production is -server option require")

func main() {
	flag.Parse()
	conf.NewConfig(*runServer)

	e := echo.New()
	connection := conf.NewDBConnection() // TODO 接続のたびに取得するほうがいいのでは？
	defer func() {
		if err := connection.Close(); err != nil {
			e.Logger.Fatal(fmt.Sprintf("Failed to close: %v", err))
		}
	}()
	interactor := interactor.NewInteractor(connection) // TODO DI用のライブラリはないのかな？
	handler := interactor.NewAppHandler()

	router.SetRoutes(e, handler)
	middleware.SetMiddlewares(e)
	if err := e.Start(fmt.Sprintf(":%d", conf.Current.Server.Port)); err != nil {
		e.Logger.Fatal(fmt.Sprintf("Failed to start: %v", err))
	}
}
