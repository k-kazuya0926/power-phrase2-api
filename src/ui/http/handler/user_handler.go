package handler

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/ui/http/request"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
	"github.com/labstack/echo"
	"github.com/olahol/go-imageupload"
)

type (
	// UserHandler interface
	UserHandler interface {
		UploadImageFile(c echo.Context) error
		CreateUser(c echo.Context) error
		Login(c echo.Context) error
		GetUser(c echo.Context) error
		UpdateUser(c echo.Context) error
		DeleteUser(c echo.Context) error
	}

	userHandler struct {
		UserUseCase usecase.UserUseCase
	}
)

// NewUserHandler UserHandlerを取得します.
func NewUserHandler(usecase usecase.UserUseCase) UserHandler {
	return &userHandler{usecase}
}

func (handler *userHandler) UploadImageFile(c echo.Context) error {
	img, err := imageupload.Process(c.Request(), "ImageFile")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	h := sha1.Sum(thumb.Data)
	fileName := fmt.Sprintf("%s_%x.png", time.Now().Format("20060102150405"), h[:])
	thumb.Save("assets/images/" + fileName)

	return c.String(http.StatusOK, "images/"+fileName)
}

func (handler *userHandler) CreateUser(c echo.Context) error {
	request := new(request.CreateUserRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err := handler.UserUseCase.CreateUser(
		request.Name,
		request.Email,
		request.Password,
		request.ImageFilePath,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (handler *userHandler) Login(c echo.Context) error {
	request := new(request.LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	userID, token, err := handler.UserUseCase.Login(request.Email, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    userID,
		"token": token,
	})
}

func (handler *userHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := &request.GetUserRequest{ID: id}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	user, err := handler.UserUseCase.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *userHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := &request.UpdateUserRequest{ID: id}
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = handler.UserUseCase.UpdateUser(
		id,
		request.Name,
		request.Email,
		request.Password,
		request.ImageFilePath,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (handler *userHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "ID：数値で入力してください。")
	}

	request := request.DeleteUserRequest{ID: id}
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if err := handler.UserUseCase.DeleteUser(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
