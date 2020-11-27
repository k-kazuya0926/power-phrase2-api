// Package handler UI層
package handler

import (
	"net/http"
	"strconv"

	"github.com/k-kazuya0926/power-phrase2-api/ui/http/request"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
	"github.com/labstack/echo"
)

type (
	// CommentHandler interface
	CommentHandler interface {
		CreateComment(c echo.Context) error
		GetComments(c echo.Context) error
		DeleteComment(c echo.Context) error
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
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}
	request := &request.CreateCommentRequest{PostID: postID}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = handler.CommentUseCase.CreateComment(
		request.PostID,
		request.UserID,
		request.Body,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// GetComments 一覧取得
func (handler *commentHandler) GetComments(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "id：数値で入力してください。")
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "limit：数値で入力してください。")
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "page：数値で入力してください。")
	}

	request := &request.GetCommentsRequest{
		PostID: postID,
		Limit:  limit,
		Page:   page,
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	totalCount, comments, err := handler.CommentUseCase.GetComments(postID, limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalCount": totalCount,
		"comments":   comments,
	})
}

// DeleteComment 削除
func (handler *commentHandler) DeleteComment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := request.DeleteCommentRequest{ID: id}
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := handler.CommentUseCase.DeleteComment(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
