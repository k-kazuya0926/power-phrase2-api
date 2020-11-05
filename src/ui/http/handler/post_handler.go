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

	postHandler struct {
		PostUseCase usecase.PostUseCase
	}
)

// NewPostHandler PostHandlerを取得します.
func NewPostHandler(usecase usecase.PostUseCase) PostHandler {
	return &postHandler{usecase}
}

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

func (handler *postHandler) GetPosts(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "limit：数値で入力してください。")
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "page：数値で入力してください。")
	}
	userID, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		userID = 0
	}

	keyword := c.QueryParam("keyword")

	request := &request.GetPostsRequest{
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
		UserID:  userID,
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	totalCount, posts, err := handler.PostUseCase.GetPosts(limit, page, keyword, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"totalCount": totalCount,
		"posts":      posts,
	})
}

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
