package usecase

import (
	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase interfase
type UserUseCase interface {
	CreateUser(user *model.User) (*model.User, error)
	Login(email, password string) (int, string, error)
	GetUsers() ([]*model.User, error)
	GetUser(id int) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id int) error
}

type userUseCase struct {
	repository.UserRepository
}

// NewUserUseCase UserUseCaseを取得します.
func NewUserUseCase(repository repository.UserRepository) UserUseCase {
	return &userUseCase{repository}
}

func (usecase *userUseCase) CreateUser(user *model.User) (*model.User, error) {
	// パスワード暗号化
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(passwordHash)

	return usecase.UserRepository.Create(user)
}

func (usecase *userUseCase) Login(email, password string) (userID int, token string, err error) {
	// TODO 整理
	var errorList []string
	if email == "" {
		errorList = append(errorList, "Emailは必須です。")
	}
	if password == "" {
		errorList = append(errorList, "パスワードは必須です。")
	}

	if len(errorList) > 0 {
		return -1, "", errors.New("バリデーションエラー") // TODO 戻り値修正
	}

	user, err := usecase.UserRepository.FetchByEmail(email)
	if err != nil {
		return -1, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, "", err
	}

	// JWTトークン生成
	token, err = createToken(user)
	if err != nil {
		return -1, "", err
	}

	return user.ID, token, err
}

func createToken(user *model.User) (string, error) {
	var err error

	// 鍵となる文字列
	// secret := os.Getenv("SECRET_KEY")
	secret := "secret" // TODO 変更

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // TODO 見直し
		"email": user.Email,
		"iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func (usecase *userUseCase) GetUsers() ([]*model.User, error) {
	return usecase.UserRepository.Fetch()
}

func (usecase *userUseCase) GetUser(id int) (*model.User, error) {
	return usecase.UserRepository.FetchByID(id)
}

func (usecase *userUseCase) UpdateUser(user *model.User) (*model.User, error) {
	return usecase.UserRepository.Update(user)
}

func (usecase *userUseCase) DeleteUser(id int) error {
	return usecase.UserRepository.Delete(id)
}
