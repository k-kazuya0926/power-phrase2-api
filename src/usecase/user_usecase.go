package usecase

import (
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
	return usecase.UserRepository.Create(user)
}

func (usecase *userUseCase) Login(email, password string) (int, string, error) {
	// TODO バリデーション

	user, err := usecase.UserRepository.FetchByEmail(email)
	// TODO エラー処理

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, "", err
	}

	// JWTトークン生成
	token, err := createToken(email)
	if err != nil {
		return -1, "", err
	}

	return user.ID, token, err
}

func createToken(email string) (string, error) {
	var err error

	// 鍵となる文字列
	// secret := os.Getenv("SECRET_KEY")
	secret := "secret" // TODO 変更

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
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
