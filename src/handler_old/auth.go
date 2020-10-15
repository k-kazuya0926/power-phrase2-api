package handler

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type jwtCustomClaims struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte(getSigningKey())

func getSigningKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env読み込みエラー")
	}

	return os.Getenv("JWT_SIGNING_KEY")
}

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

func Signup(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	if err := validateForSignup(user); err != nil {
		return err
	}

	// パスワード暗号化
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)

	model.CreateUser(user)
	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

func validateForSignup(user *model.User) *echo.HTTPError {
	// ユーザー名またはパスワードが空である場合、エラー
	if user.Name == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		}
	}

	// ユーザー名が登録済みである場合、エラー
	if u := model.FindUser(&model.User{Name: user.Name}); u.ID != 0 {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "name already exists",
		}
	}

	return nil
}

func Login(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	// ユーザー名で検索
	user := model.FindUser(&model.User{Name: u.Name})
	// ユーザーが存在しない、またはパスワードが違う場合エラー
	if user.ID == 0 || user.Password != u.Password {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "invalid name or password",
		}
	}

	claims := &jwtCustomClaims{
		user.ID,
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func userIDFromToken(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}