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
	// PostHandler interface
	PostHandler interface {
		CreatePost(c echo.Context) error
		GetPosts(c echo.Context) error
		GetPost(c echo.Context) error
		UpdatePost(c echo.Context) error
		DeletePost(c echo.Context) error
	}

	// postHandler 構造体
	postHandler struct {
		PostUseCase usecase.PostUseCase
	}
)

// NewPostHandler PostHandlerを生成。
func NewPostHandler(usecase usecase.PostUseCase) PostHandler {
	return &postHandler{usecase}
}

// CreatePost 登録
func (handler *postHandler) CreatePost(c echo.Context) error {
	request := new(request.CreatePostRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err := handler.PostUseCase.CreatePost(
		request.UserID,
		request.Title,
		request.Speaker,
		request.Detail,
		request.MovieURL,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// GetPosts 一覧取得
func (handler *postHandler) GetPosts(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "limit：数値で入力してください。")
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "page：数値で入力してください。")
	}
	postUserID, err := strconv.Atoi(c.QueryParam("post_user_id"))
	if err != nil {
		postUserID = 0
	}
	loginUserID, err := strconv.Atoi(c.QueryParam("login_user_id"))
	if err != nil {
		loginUserID = 0
	}

	keyword := c.QueryParam("keyword")

	request := &request.GetPostsRequest{
		Limit:       limit,
		Page:        page,
		Keyword:     keyword,
		PostUserID:  postUserID,
		LoginUserID: loginUserID,
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	totalCount, posts, err := handler.PostUseCase.GetPosts(limit, page, keyword, postUserID, loginUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalCount": totalCount,
		"posts":      posts,
	})
}

// GetPost　1件取得
func (handler *postHandler) GetPost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := &request.GetPostRequest{ID: id}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	post, err := handler.PostUseCase.GetPost(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, post)
}

// UpdatePost 更新
func (handler *postHandler) UpdatePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := &request.UpdatePostRequest{ID: id}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = handler.PostUseCase.UpdatePost(
		id,
		request.Title,
		request.Speaker,
		request.Detail,
		request.MovieURL,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DeletePost 削除
func (handler *postHandler) DeletePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := request.DeletePostRequest{ID: id}
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := handler.PostUseCase.DeletePost(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
