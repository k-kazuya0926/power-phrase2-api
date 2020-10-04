package model

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type Post struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    uint       `json:"user_id"`
	Title     string     `json:"title"`
	Speaker   string     `json:"speaker"`
	Detail    string     `json:"detail"`
	MovieURL  string     `json:"movie_url"`
}

func GetAllPosts(c echo.Context) error {
	db := Connection()
	defer db.Close()
	db.AutoMigrate(&Post{})

	var posts []Post
	db.Find(&posts)
	return c.JSON(http.StatusOK, posts)
}

func GetPost(c echo.Context) error {
	db := Connection()
	defer db.Close()
	db.AutoMigrate(&Post{})

	if id := c.Param("id"); id != "" {
		var post Post
		db.First(&post, id)
		return c.JSON(http.StatusOK, post)
	} else {
		return c.JSON(http.StatusNotFound, nil)
	}
}

func CreatePost(c echo.Context) error {
	db := Connection()
	defer db.Close()
	db.AutoMigrate(&Post{})

	post := new(Post)
	if err := c.Bind(post); err != nil {
		return err
	}
	db.Create(&post)

	return c.JSON(http.StatusOK, post)
}

func UpdatePost(c echo.Context) error {
	db := Connection()
	defer db.Close()

	newPost := new(Post)
	if err := c.Bind(newPost); err != nil {
		return err
	}

	if id := c.Param("id"); id != "" {
		var post Post
		db.First(&post, id).Update(newPost)
		return c.JSON(http.StatusOK, post)
	} else {
		return c.JSON(http.StatusNotFound, nil)
	}

}

func DeletePost(c echo.Context) error {
	db := Connection()
	defer db.Close()

	if id := c.Param("id"); id != "" {
		var post Post
		db.First(&post, id)
		db.Delete(post)
		return c.JSON(http.StatusOK, post)
	} else {
		return c.JSON(http.StatusNotFound, nil)
	}
}
