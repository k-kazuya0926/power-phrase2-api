package datastore

import (
	"fmt"
	"testing"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestPostRepository_Create(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	repository := &postRepository{}
	postForInput := getMockPost(userForInput.ID)

	// 2. Exercise
	err := repository.Create(postForInput)

	// 3. Verify
	assert.NoError(t, err)

	// 件数
	var count int
	db.Table("posts").Count(&count)
	assert.Equal(t, 1, count)

	// 内容
	post := model.Post{}
	db.First(&post)
	assert.Equal(t, postForInput.UserID, post.UserID)
	assert.Equal(t, postForInput.Title, post.Title)
	assert.Equal(t, postForInput.Speaker, post.Speaker)
	assert.Equal(t, postForInput.Detail, post.Detail)
	assert.Equal(t, postForInput.MovieURL, post.MovieURL)

	// 4. Teardown
	teardown(db)
}

func TestPostRepository_Fetch(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	postForInput2 := getMockPost(userForInput.ID)
	db.Create(&postForInput2)

	repository := &postRepository{}

	// 2. Exercise
	posts, err := repository.Fetch(1, 1, "")

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, postForInput2.UserID, posts[0].UserID)
	assert.Equal(t, postForInput2.Speaker, posts[0].Speaker)
	assert.Equal(t, postForInput2.Detail, posts[0].Detail)
	assert.Equal(t, postForInput2.MovieURL, posts[0].MovieURL)

	// 4. Teardown
	teardown(db)
}

func TestPostRepository_FetchById(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	repository := &postRepository{}

	// 2. Exercise
	actualPost, err := repository.FetchByID(postForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, postForInput.ID, actualPost.ID)
	assert.Equal(t, postForInput.Title, actualPost.Title)
	assert.Equal(t, postForInput.Speaker, actualPost.Speaker)
	assert.Equal(t, postForInput.Detail, actualPost.Detail)
	assert.Equal(t, postForInput.MovieURL, actualPost.MovieURL)

	// 4. Teardown
	teardown(db)
}

func TestPostRepository_Update(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)
	postForInput.Title = "title2"
	postForInput.Speaker = "speaker2"
	postForInput.Detail = "detail2"
	postForInput.MovieURL = "http://www.example.com/2"

	repository := &postRepository{}

	// 2. Exercise
	err := repository.Update(postForInput)

	// 3. Verify
	assert.NoError(t, err)

	// 件数
	var count int
	db.Table("posts").Count(&count)
	assert.Equal(t, 1, count)

	// 内容
	post := model.Post{}
	db.First(&post)
	assert.Equal(t, postForInput.UserID, post.UserID)
	assert.Equal(t, postForInput.Title, post.Title)
	assert.Equal(t, postForInput.Speaker, post.Speaker)
	assert.Equal(t, postForInput.Detail, post.Detail)
	assert.Equal(t, postForInput.MovieURL, post.MovieURL)

	// 4. Teardown
	teardown(db)
}

func TestPostRepository_Delete(t *testing.T) {
	// 1. Setup
	db := conf.NewTestDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	repository := &postRepository{}

	// 2. Exercise
	err := repository.Delete(postForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	assert.Error(t, db.First(&postForInput).Error)

	// 4. Teardown
	teardown(db)
}

func getMockPost(userID int) *model.Post {
	post := &model.Post{
		UserID:   userID,
		Title:    fmt.Sprintf("title%d", userID),
		Speaker:  fmt.Sprintf("speaker%d", userID),
		Detail:   fmt.Sprintf("detail%d", userID),
		MovieURL: fmt.Sprintf("http://www.example.com/%d", userID),
	}
	return post
}
