package usecase

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase interfase
type UserUseCase interface {
	CreateUser(name, email, password, imageURL string) (userID int, err error)
	Login(email, password string) (token string, err error)
	GetUser(id int) (*model.User, error)
	UpdateUser(userID int, name, email, password, imageURL string) error
	DeleteUser(id int) error
}

type userUseCase struct {
	repository.UserRepository
}

// NewUserUseCase UserUseCaseを取得します.
func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}

func (usecase *userUseCase) CreateUser(name, email, password, imageURL string) (userID int, err error) {
	// パスワード暗号化
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(passwordHash),
		ImageURL: imageURL,
	}
	err = usecase.UserRepository.Create(&user)

	return user.ID, err
}

func (usecase *userUseCase) Login(email, password string) (token string, err error) {
	user, err := usecase.UserRepository.FetchByEmail(email)
	if err != nil {
		return "", errors.New("ユーザーの取得に失敗しました。")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("メールアドレスまたはパスワードに誤りがあります。")
	}

	// JWTトークン生成
	token, err = createToken(user)

	return token, err
}

func createToken(user *model.User) (string, error) {
	// 鍵となる文字列
	// secret := os.Getenv("SECRET_KEY")
	secret := "secret" // TODO 変更

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID
	claims["name"] = user.Name
	claims["picture"] = "TODO" // TODO
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))

	return t, err
}

func (usecase *userUseCase) GetUser(id int) (*model.User, error) {
	user, err := usecase.UserRepository.FetchByID(id)
	if err != nil {
		return nil, errors.New("ユーザーの取得に失敗しました。")
	}
	return user, nil
}

func (usecase *userUseCase) UpdateUser(userID int, name, email, password, imageURL string) error {
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
		ID:       userID,
		Name:     name,
		Email:    email,
		Password: newPassword,
		ImageURL: imageURL,
	}
	if err := usecase.UserRepository.Update(&user); err != nil {
		return errors.New("ユーザーの更新に失敗しました。")
	}
	return nil
}

func (usecase *userUseCase) DeleteUser(id int) error {
	if err := usecase.UserRepository.Delete(id); err != nil {
		return errors.New("ユーザーの削除に失敗しました。")
	}
	return nil
}
