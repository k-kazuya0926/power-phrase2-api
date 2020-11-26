// Package handler UI層
package handler

import (
	"net/http"

	"github.com/k-kazuya0926/power-phrase2-api/ui/http/request"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
	"github.com/labstack/echo"
)

type (
	// CommentHandler interface
	CommentHandler interface {
		CreateComment(c echo.Context) error
	}

	// commentHandler 構造体
	commentHandler struct {
		CommentUseCase usecase.CommentUseCase
	}
)

// NewCommentHandler CommentHandlerを生成。
func NewCommentHandler(usecase usecase.CommentUseCase) CommentHandler {
	return &commentHandler{usecase}
}

// CreateComment 登録
func (handler *commentHandler) CreateComment(c echo.Context) error {
	request := new(request.CreateCommentRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err := handler.CommentUseCase.CreateComment(
		request.PostID,
		request.UserID,
		request.Body,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
