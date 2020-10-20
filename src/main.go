package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/interactor"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/middleware"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/router"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err == nil {
		return err
	}

	var errorMessages []string //バリデーションでNGとなった独自エラーメッセージを格納
	for _, err := range err.(validator.ValidationErrors) {

		// TODO 整理
		var errorMessage string
		fieldName := err.Field() //バリデーションでNGになった変数名を取得

		switch fieldName {
		case "Name":
			errorMessage = "error message for Name"
		case "Password":
			errorMessage = "error message for Password"
		case "Email":
			var typ = err.Tag() //バリデーションでNGになったタグ名を取得
			switch typ {
			case "required":
				errorMessage = "emailは必須です。"
			case "email":
				errorMessage = "emailの形式が正しくありません。"
			}
		}
		errorMessages = append(errorMessages, errorMessage)
	}
	return errors.New(strings.Join(errorMessages, "\n"))
}

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

	e.Validator = &CustomValidator{validator: validator.New()}

	if err := e.Start(fmt.Sprintf(":%d", conf.Current.Server.Port)); err != nil {
		e.Logger.Fatal(fmt.Sprintf("Failed to start: %v", err))
	}
}
