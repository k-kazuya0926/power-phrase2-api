package datastore

import (
	"fmt"
	"testing"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
)

// 投稿登録
func TestPostRepository_Create(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	repository := &postRepository{}
	postForInput := makePost(userForInput.ID)

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

// 投稿一覧取得
func TestPostRepository_Fetch(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	postForInput2 := makePost(userForInput.ID)
	db.Create(&postForInput2)

	repository := &postRepository{}

	postUserID := 0  // TODO 投稿ユーザーID指定がある場合
	loginUserID := 0 // TODO ログインユーザーID指定がある場合

	// 2. Exercise
	totalCount, posts, err := repository.Fetch(1, 1, "", postUserID, loginUserID)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, 2, totalCount)
	assert.Equal(t, postForInput2.UserID, posts[0].UserID)
	assert.Equal(t, postForInput2.Speaker, posts[0].Speaker)
	assert.Equal(t, postForInput2.Detail, posts[0].Detail)
	assert.Equal(t, postForInput2.MovieURL, posts[0].MovieURL)
	assert.Equal(t, userForInput.Name, posts[0].UserName)

	// 4. Teardown
	teardown(db)
}

// 投稿詳細取得
func TestPostRepository_FetchById(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	postForInput2 := makePost(userForInput.ID)
	db.Create(&postForInput2)
	db.First(&postForInput2)

	repository := &postRepository{}

	// 2. Exercise
	actualPost, err := repository.FetchByID(postForInput.ID, userForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, postForInput.ID, actualPost.ID)
	assert.Equal(t, postForInput.Title, actualPost.Title)
	assert.Equal(t, postForInput.Speaker, actualPost.Speaker)
	assert.Equal(t, postForInput.Detail, actualPost.Detail)
	assert.Equal(t, postForInput.MovieURL, actualPost.MovieURL)
	assert.Equal(t, userForInput.Name, actualPost.UserName)

	// 4. Teardown
	teardown(db)
}

// 投稿更新
func TestPostRepository_Update(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)
	postForInput.Title = "title2"
	postForInput.Speaker = "speaker2"
	postForInput.Detail = "detail2"
	postForInput.MovieURL = "https://www.example.com/watch?v=2"

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

// 投稿削除
func TestPostRepository_Delete(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
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

// コメント登録
func TestPostRepository_CreateComment(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	repository := &postRepository{}
	commentForInput := makeComment(postForInput.ID, userForInput.ID)

	// 2. Exercise
	err := repository.CreateComment(commentForInput)

	// 3. Verify
	assert.NoError(t, err)

	// 件数
	var count int
	db.Table("comments").Count(&count)
	assert.Equal(t, 1, count)

	// 内容
	comment := model.Comment{}
	db.First(&comment)
	assert.Equal(t, commentForInput.PostID, comment.PostID)
	assert.Equal(t, commentForInput.UserID, comment.UserID)
	assert.Equal(t, commentForInput.Body, comment.Body)

	// 4. Teardown
	teardown(db)
}

// コメント一覧取得
func TestPostRepository_FetchComments(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	commentForInput := makeComment(postForInput.ID, userForInput.ID)
	db.Create(&commentForInput)
	commentForInput2 := makeComment(postForInput.ID, userForInput.ID)
	db.Create(&commentForInput2)

	repository := &postRepository{}

	// 2. Exercise
	totalCount, comments, err := repository.FetchComments(postForInput.ID, 1, 1)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, 2, totalCount)
	assert.Equal(t, commentForInput2.PostID, comments[0].PostID)
	assert.Equal(t, commentForInput2.UserID, comments[0].UserID)
	assert.Equal(t, commentForInput2.Body, comments[0].Body)
	assert.Equal(t, userForInput.Name, comments[0].UserName)
	assert.Equal(t, userForInput.ImageFilePath, comments[0].UserImageFilePath)

	// 4. Teardown
	teardown(db)
}

// コメント削除
func TestPostRepository_DeleteComment(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	commentForInput := makeComment(postForInput.ID, userForInput.ID)
	db.Create(&commentForInput)
	db.First(&commentForInput)

	repository := &postRepository{}

	// 2. Exercise
	err := repository.DeleteComment(commentForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	assert.Error(t, db.First(&commentForInput).Error)

	// 4. Teardown
	teardown(db)
}

// お気に入り登録
func TestPostRepository_CreateFavorite(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	// ユーザー
	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	// 投稿
	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	// お気に入り
	favoriteForInput := makeFavorite(userForInput.ID, postForInput.ID)

	repository := &postRepository{}

	// 2. Exercise
	err := repository.CreateFavorite(favoriteForInput)

	// 3. Verify
	assert.NoError(t, err)

	// 件数
	var count int
	db.Table("favorites").Count(&count)
	assert.Equal(t, 1, count)

	// 内容
	favorite := model.Favorite{}
	db.First(&favorite)
	assert.Equal(t, favoriteForInput.UserID, favorite.UserID)
	assert.Equal(t, favoriteForInput.PostID, favorite.PostID)

	// 4. Teardown
	teardown(db)
}

// お気に入り一覧取得
func TestPostRepository_FetchFavorites(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	// ユーザー
	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	// 投稿
	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	postForInput2 := makePost(userForInput.ID)
	db.Create(&postForInput2)
	db.First(&postForInput2)

	// お気に入り
	favoriteForInput := makeFavorite(userForInput.ID, postForInput.ID)
	db.Create(&favoriteForInput)

	favoriteForInput2 := makeFavorite(userForInput.ID, postForInput2.ID)
	db.Create(&favoriteForInput2)

	repository := &postRepository{}

	// 2. Exercise
	totalCount, favorites, err := repository.FetchFavorites(userForInput.ID, 10, 1)

	// 3. Verify
	assert.NoError(t, err)

	// 内容
	assert.Equal(t, 2, totalCount)
	assert.Equal(t, favoriteForInput2.PostID, favorites[0].ID)
	assert.Equal(t, favoriteForInput2.UserID, favorites[0].UserID)
	assert.Equal(t, userForInput.Name, favorites[0].UserName)
	assert.Equal(t, userForInput.ImageFilePath, favorites[0].UserImageFilePath)

	// 4. Teardown
	teardown(db)
}

// お気に入り削除
func TestPostRepository_DeleteFavorite(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	// ユーザー
	userForInput := makeUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	// 投稿
	postForInput := makePost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	// お気に入り
	favoriteForInput := makeFavorite(userForInput.ID, postForInput.ID)
	db.Create(&favoriteForInput)
	db.First(&favoriteForInput)

	repository := &postRepository{}

	// 2. Exercise
	err := repository.DeleteFavorite(userForInput.ID, postForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	assert.Error(t, db.First(&favoriteForInput).Error)

	// 4. Teardown
	teardown(db)
}

// Postを生成
func makePost(userID int) *model.Post {
	return &model.Post{
		UserID:   userID,
		Title:    fmt.Sprintf("title%d", userID),
		Speaker:  fmt.Sprintf("speaker%d", userID),
		Detail:   fmt.Sprintf("detail%d", userID),
		MovieURL: fmt.Sprintf("https://www.example.com/watch?v=%d", userID),
	}
}

// GetPostResultを生成
func makeGetPostResult(userID int) *model.GetPostResult {
	return &model.GetPostResult{
		Post:              *makePost(userID),
		UserName:          fmt.Sprintf("testuser%d", userID),
		UserImageFilePath: fmt.Sprintf("images/%d.png", userID),
		CommentCount:      userID,
	}
}

// Commentを生成
func makeComment(postID, userID int) *model.Comment {
	return &model.Comment{
		PostID: postID,
		UserID: userID,
		Body:   "body",
	}
}

// Favoriteを生成
func makeFavorite(userID, postID int) *model.Favorite {
	return &model.Favorite{
		UserID: userID,
		PostID: postID,
	}
}
