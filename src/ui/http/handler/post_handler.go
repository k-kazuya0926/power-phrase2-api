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
		// 投稿登録
		CreatePost(c echo.Context) error
		// 投稿一覧取得
		GetPosts(c echo.Context) error
		// 投稿詳細取得
		GetPost(c echo.Context) error
		// 投稿更新
		UpdatePost(c echo.Context) error
		// 投稿削除
		DeletePost(c echo.Context) error

		// お気に入り登録
		CreateFavorite(c echo.Context) error
		// お気に入り一覧取得
		GetFavorites(c echo.Context) error
		// お気に入り削除
		DeleteFavorite(c echo.Context) error
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

// CreatePost 投稿登録
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

// GetPosts 投稿一覧取得
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

// GetPost　投稿詳細取得
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

// UpdatePost 投稿更新
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

// DeletePost 投稿削除
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

// CreateFavorite お気に入り登録
func (handler *postHandler) CreateFavorite(c echo.Context) error {
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

	err = handler.PostUseCase.CreateFavorite(
		request.UserID,
		request.PostID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// GetFavorites お気に入り一覧取得
func (handler *postHandler) GetFavorites(c echo.Context) error {
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

	totalCount, posts, err := handler.PostUseCase.GetFavorites(userID, limit, page)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalCount": totalCount,
		"posts":      posts,
	})
}

// DeleteFavorite お気に入り削除
func (handler *postHandler) DeleteFavorite(c echo.Context) error {
	request := new(request.DeleteFavoriteRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := handler.PostUseCase.DeleteFavorite(request.UserID, request.PostID); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
