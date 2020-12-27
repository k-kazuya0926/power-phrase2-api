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
	// FavoriteHandler interface
	FavoriteHandler interface {
		CreateFavorite(c echo.Context) error
		GetFavorites(c echo.Context) error
		DeleteFavorite(c echo.Context) error
	}

	// favoriteHandler 構造体
	favoriteHandler struct {
		FavoriteUseCase usecase.FavoriteUseCase
	}
)

// NewFavoriteHandler FavoriteHandlerを生成。
func NewFavoriteHandler(usecase usecase.FavoriteUseCase) FavoriteHandler {
	return &favoriteHandler{usecase}
}

// CreateFavorite お気に入り登録
func (handler *favoriteHandler) CreateFavorite(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}
	request := &request.CreateFavoriteRequest{PostID: postID}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = handler.FavoriteUseCase.CreateFavorite(
		request.PostID,
		request.UserID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// GetFavorites お気に入り一覧取得
func (handler *favoriteHandler) GetFavorites(c echo.Context) error {
	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "user_id：数値で入力してください。")
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "limit：数値で入力してください。")
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "page：数値で入力してください。")
	}

	request := &request.GetFavoritesRequest{
		UserID: userID,
		Limit:  limit,
		Page:   page,
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	totalCount, posts, err := handler.FavoriteUseCase.GetFavorites(userID, limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalCount": totalCount,
		"posts":  posts,
	})
}

// DeleteFavorite お気に入り削除
func (handler *favoriteHandler) DeleteFavorite(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := request.DeleteFavoriteRequest{ID: id}
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := handler.FavoriteUseCase.DeleteFavorite(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
