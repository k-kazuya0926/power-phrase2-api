package model

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	ImageURL  string     `json:"image_url"`
}

func CreateUser(user *User) {
	db := Connection()
	defer db.Close()
	db.AutoMigrate(&User{})
	db.Create(user)
}

func FindUser(u *User) User {
	var user User
	db := Connection()
	db.AutoMigrate(&User{})
	db.Where(u).First(&user)
	return user
}

// TODO 以下、echoに依存しないように見直し
func GetAllUsers(c echo.Context) error {
	db := Connection()
	defer db.Close()
	db.AutoMigrate(&User{})

	var users []User
	db.Find(&users)
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	db := Connection()
	defer db.Close()
	db.AutoMigrate(&User{})

	if id := c.Param("id"); id != "" {
		var user User
		db.First(&user, id)
		return c.JSON(http.StatusOK, user)
	} else {
		return c.JSON(http.StatusNotFound, nil)
	}
}

func UpdateUser(c echo.Context) error {
	db := Connection()
	defer db.Close()

	newUser := new(User)
	if err := c.Bind(newUser); err != nil {
		return err
	}

	if id := c.Param("id"); id != "" {
		var user User
		db.First(&user, id).Update(newUser)
		return c.JSON(http.StatusOK, user)
	} else {
		return c.JSON(http.StatusNotFound, nil)
	}

}

func DeleteUser(c echo.Context) error {
	db := Connection()
	defer db.Close()

	if id := c.Param("id"); id != "" {
		var user User
		db.First(&user, id)
		db.Delete(user)
		return c.JSON(http.StatusOK, user)
	} else {
		return c.JSON(http.StatusNotFound, nil)
	}
}
