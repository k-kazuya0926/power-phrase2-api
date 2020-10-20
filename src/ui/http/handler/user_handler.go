package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/ui/http/request"
	"github.com/k-kazuya0926/power-phrase2-api/usecase"
	"github.com/labstack/echo"
)

type (
	// UserHandler interface
	UserHandler interface {
		CreateUser(c echo.Context) error
		Login(c echo.Context) error
		GetUsers(c echo.Context) error
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

func (handler *userHandler) CreateUser(c echo.Context) error {
	request := new(request.CreateUserRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	if err := c.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := handler.UserUseCase.CreateUser(
		request.Name,
		request.Email,
		request.Password,
		request.ImageURL,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user.ID)
}

func (handler *userHandler) Login(c echo.Context) error {
	userID, token, err := handler.UserUseCase.Login(c.FormValue("email"), c.FormValue("password"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ログイン時にエラーが発生しました。")
	}

	return c.JSON(http.StatusCreated, fmt.Sprintf("UserID: %d, token: %s", userID, token)) // TODO 戻り値修正
}

func (handler *userHandler) GetUsers(c echo.Context) error {
	users, err := handler.UserUseCase.GetUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Users does not exist.")
	}

	return c.JSON(http.StatusOK, users)
}

func (handler *userHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID must be int")
	}

	user, err := handler.UserUseCase.GetUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User does not exist.")
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *userHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID must be int")
	}

	newUser := new(model.User)
	if err := c.Bind(newUser); err != nil {
		return err
	}
	newUser.ID = id

	user, err := handler.UserUseCase.UpdateUser(newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "User can not Create.")
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *userHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID must be int")
	}

	if err := handler.UserUseCase.DeleteUser(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "User can not Delete.")
	}

	return c.NoContent(http.StatusNoContent)
}
