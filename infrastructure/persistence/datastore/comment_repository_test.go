package datastore

import (
	"testing"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestCommentRepository_Create(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	repository := &commentRepository{}
	commentForInput := getMockComment(postForInput.ID, userForInput.ID)

	// 2. Exercise
	err := repository.Create(commentForInput)

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

func TestCommentRepository_Fetch(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	commentForInput := getMockComment(postForInput.ID, userForInput.ID)
	db.Create(&commentForInput)
	commentForInput2 := getMockComment(postForInput.ID, userForInput.ID)
	db.Create(&commentForInput2)

	repository := &commentRepository{}

	// 2. Exercise
	totalCount, comments, err := repository.Fetch(postForInput.ID, 1, 1)

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

func TestCommentRepository_Delete(t *testing.T) {
	// 1. Setup
	setup()
	db := conf.NewDBConnection()
	defer db.Close()

	userForInput := getMockUserForInput(1)
	db.Create(&userForInput)
	db.First(&userForInput)

	postForInput := getMockPost(userForInput.ID)
	db.Create(&postForInput)
	db.First(&postForInput)

	commentForInput := getMockComment(postForInput.ID, userForInput.ID)
	db.Create(&commentForInput)
	db.First(&commentForInput)

	repository := &commentRepository{}

	// 2. Exercise
	err := repository.Delete(commentForInput.ID)

	// 3. Verify
	assert.NoError(t, err)

	assert.Error(t, db.First(&commentForInput).Error)

	// 4. Teardown
	teardown(db)
}

func getMockComment(postID, userID int) *model.Comment {
	return &model.Comment{
		PostID: postID,
		UserID: userID,
		Body:   "body",
	}
}
