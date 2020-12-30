package handler

import (
	"io"
	"net/http/httptest"

	"github.com/k-kazuya0926/power-phrase2-api/validator"
	"github.com/labstack/echo"
)

func createContext(method, path string, body io.Reader, rec *httptest.ResponseRecorder) echo.Context {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	e := echo.New()
	e.Validator = validator.NewValidator()
	return e.NewContext(req, rec)
}
