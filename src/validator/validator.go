package validator

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err == nil {
		return err
	}

	var errorMessages []string //バリデーションでNGとなった独自エラーメッセージを格納
	for _, err := range err.(validator.ValidationErrors) {
		var errorMessage string

		var typ = err.Tag()
		switch typ {
		case "required":
			errorMessage = fmt.Sprintf("%s：必須です。", err.Field())
		case "email":
			errorMessage = fmt.Sprintf("%s：正しい形式で入力してください。", err.Field())
		case "min":
			errorMessage = fmt.Sprintf("%s：%s以上の値を入力してください。", err.Field(), err.Param())
		default:
			errorMessage = fmt.Sprintf("%s：正しい値を入力してください。", err.Field())
		}

		errorMessages = append(errorMessages, errorMessage)
	}
	return errors.New(strings.Join(errorMessages, "\n"))
}
