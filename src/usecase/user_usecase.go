// Package usecase Application Service層。
package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase インターフェース
type UserUseCase interface {
	CreateUser(name, email, password, imageFilePath string) (userID int, token string, err error)
	Login(email, password string) (userID int, token string, err error)
	GetUser(id int) (*model.User, error)
	UpdateUser(userID int, name, email, password, imageFilePath string) error
	DeleteUser(id int) error
}

// userUseCase 構造体
type userUseCase struct {
	repository.UserRepository
}

// NewUserUseCase UserUseCaseを生成。
func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}

// CreateUser 登録
func (usecase *userUseCase) CreateUser(name, email, password, imageFilePath string) (userID int, token string, err error) {
	// パスワード暗号化
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}

	user := model.User{
		Name:          name,
		Email:         email,
		Password:      string(passwordHash),
		ImageFilePath: imageFilePath,
	}

	if err = usecase.UserRepository.Create(&user); err != nil {
		return 0, "", err
	}

	// JWTトークン生成
	token, err = createToken(&user)

	return user.ID, token, err
}

// Login ログイン
func (usecase *userUseCase) Login(email, password string) (userID int, token string, err error) {
	user, err := usecase.UserRepository.FetchByEmail(email)
	if err != nil {
		return 0, "", errors.New("メールアドレスまたはパスワードに誤りがあります。")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return 0, "", errors.New("メールアドレスまたはパスワードに誤りがあります。")
	}

	// JWTトークン生成
	token, err = createToken(user)

	return user.ID, token, err
}

// createToken JWTトークンを生成
func createToken(user *model.User) (string, error) {
	// 鍵となる文字列
	secret := os.Getenv("JWT_SIGNING_KEY")

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))

	return t, err
}

// GetUser 1件取得
func (usecase *userUseCase) GetUser(id int) (*model.User, error) {
	user, err := usecase.UserRepository.FetchByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新
func (usecase *userUseCase) UpdateUser(userID int, name, email, password, imageFilePath string) error {
	oldUser, err := usecase.UserRepository.FetchByID(userID)
	if err != nil {
		return err
	}

	// パスワード暗号化
	newPassword := ""
	if password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		newPassword = string(passwordHash)
	}

	newUser := model.User{
		ID:            userID,
		Name:          name,
		Email:         email,
		Password:      newPassword,
		ImageFilePath: imageFilePath,
	}
	if err := usecase.UserRepository.Update(&newUser); err != nil {
		return err
	}

	if imageFilePath != "" {
		os.Remove("assets/" + oldUser.ImageFilePath)
	}

	return nil
}

// DeleteUser 削除
func (usecase *userUseCase) DeleteUser(id int) error {
	if err := usecase.UserRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
