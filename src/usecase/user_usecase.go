package usecase

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	CreateUser(name, email, password, imageFilePath string) (userID int, token string, err error)
	Login(email, password string) (userID int, token string, err error)
	GetUser(id int) (*model.User, error)
	UpdateUser(userID int, name, email, password, imageFilePath string) error
	DeleteUser(id int) error
}

type userUseCase struct {
	repository.UserRepository
}

// NewUserUseCase UserUseCaseを取得します.
func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}

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

func createToken(user *model.User) (string, error) {
	// 鍵となる文字列
	secret := conf.Current.Jwt.Secret

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

func (usecase *userUseCase) GetUser(id int) (*model.User, error) {
	user, err := usecase.UserRepository.FetchByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *userUseCase) UpdateUser(userID int, name, email, password, imageFilePath string) error {
	// パスワード暗号化
	newPassword := ""
	if password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		newPassword = string(passwordHash)
	}

	user := model.User{
		ID:            userID,
		Name:          name,
		Email:         email,
		Password:      newPassword,
		ImageFilePath: imageFilePath,
	}
	if err := usecase.UserRepository.Update(&user); err != nil {
		return err
	}
	return nil
}

func (usecase *userUseCase) DeleteUser(id int) error {
	if err := usecase.UserRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
