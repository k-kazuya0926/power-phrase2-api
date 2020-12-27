package usecase

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock
type mockPostRepository struct {
	mock.Mock
}

// 投稿登録
func (repository *mockPostRepository) Create(post *model.Post) error {
	return repository.Called(post).Error(0)
}

// 投稿一覧取得
func (repository *mockPostRepository) Fetch(limit, page int, keyword string, postUserID, loginUserID int) (int, []*model.GetPostResult, error) {
	args := repository.Called(limit, page, keyword, postUserID, loginUserID)
	posts, ok := args.Get(1).([]*model.GetPostResult)
	if ok {
		return args.Int(0), posts, args.Error(2)
	}

	return args.Int(0), nil, args.Error(2)
}

// 投稿詳細取得
func (repository *mockPostRepository) FetchByID(id int) (*model.GetPostResult, error) {
	args := repository.Called(id)
	post, ok := args.Get(0).(*model.GetPostResult)
	if ok {
		return post, args.Error(1)
	}

	return nil, args.Error(1)
}

// 投稿更新
func (repository *mockPostRepository) Update(post *model.Post) error {
	return repository.Called(post).Error(0)
}

// 投稿削除
func (repository *mockPostRepository) Delete(id int) error {
	return repository.Called(id).Error(0)
}

// コメント登録
func (repository *mockPostRepository) CreateComment(comment *model.Comment) error {
	return repository.Called(comment).Error(0)
}

// コメント一覧取得
func (repository *mockPostRepository) FetchComments(postID, limit, page int) (int, []*model.GetCommentResult, error) {
	args := repository.Called(postID, limit, page)
	comments, ok := args.Get(1).([]*model.GetCommentResult)
	if ok {
		return args.Int(0), comments, args.Error(2)
	}

	return args.Int(0), nil, args.Error(2)
}

// コメント削除
func (repository *mockPostRepository) DeleteComment(id int) error {
	return repository.Called(id).Error(0)
}

// お気に入り登録
func (repository *mockPostRepository) CreateFavorite(favorite *model.Favorite) error {
	return repository.Called(favorite).Error(0)
}

// お気に入り一覧取得
func (repository *mockPostRepository) FetchFavorites(userID, limit, page int) (int, []*model.GetPostResult, error) {
	args := repository.Called(userID, limit, page)
	comments, ok := args.Get(1).([]*model.GetPostResult)
	if ok {
		return args.Int(0), comments, args.Error(2)
	}

	return args.Int(0), nil, args.Error(2)
}

// お気に入り削除
func (repository *mockPostRepository) DeleteFavorite(userID, postID int) error {
	return repository.Called(userID, postID).Error(0)
}

// 入力用投稿を生成
func makePostForInput(id int) *model.Post {
	post := &model.Post{
		ID:       id,
		UserID:   id,
		Title:    fmt.Sprintf("title%d", id),
		Speaker:  fmt.Sprintf("speaker%d", id),
		Detail:   fmt.Sprintf("detail%d", id),
		MovieURL: fmt.Sprintf("https://www.example.com/watch?v=%d", id),
	}
	return post
}

// DBから取得された投稿を生成
func makePostForRead(id int) *model.Post {
	post := makePostForInput(id)
	post.CreatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	post.UpdatedAt = time.Date(2015, 9, 13, 12, 35, 42, 123456789, time.Local)
	return post
}

func makeGetPostResult(id int) *model.GetPostResult {
	post := makePostForRead(id)
	return &model.GetPostResult{
		Post:              *post,
		UserName:          fmt.Sprintf("username%d", id),
		UserImageFilePath: fmt.Sprintf("images/%d.png", id),
	}
}

// お気に入りを生成
func makeFavorite(userID, postID int) *model.Favorite {
	return &model.Favorite{
		UserID: userID,
		PostID: postID,
	}
}

// 投稿登録テスト
func TestCreatePost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := makePostForInput(id)
	repository.On("Create", mock.AnythingOfType("*model.Post")).Return(nil)

	// 2. Exercise
	err := usecase.CreatePost(post.UserID, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestCreatePost_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := makePostForInput(id)
	repository.On("Create", mock.AnythingOfType("*model.Post")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.CreatePost(post.UserID, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// 投稿一覧テスト
func TestGetPosts_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	limit := 3
	page := 1
	keyword := ""
	postUserID := 0  // TODO 投稿ユーザーID指定がある場合
	loginUserID := 0 // TODO ログインユーザーID指定がある場合
	expectedTotalCount := 2
	expectedPosts := []*model.GetPostResult{makeGetPostResult(1), makeGetPostResult(2)}
	repository.On("Fetch", limit, page, keyword, postUserID, loginUserID).Return(expectedTotalCount, expectedPosts, nil)

	// 2. Exercise
	totalCount, posts, err := usecase.GetPosts(limit, page, keyword, postUserID, loginUserID)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, expectedTotalCount, totalCount)
	assert.Equal(t, len(expectedPosts), len(posts))
	assert.Equal(t, expectedPosts[0], posts[0])
	assert.Equal(t, expectedPosts[1], posts[1])

	// 4. Teardown
}

func TestGetPosts_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	limit := 3
	page := 1
	keyword := ""
	postUserID := 0  // TODO 投稿ユーザーID指定がある場合
	loginUserID := 0 // TODO ログインユーザーID指定がある場合
	repository.On("Fetch", limit, page, keyword, postUserID, loginUserID).Return(0, nil, errors.New("error"))

	// 2. Execise
	totalCount, posts, err := usecase.GetPosts(limit, page, keyword, postUserID, loginUserID)

	// 3. Verify
	assert.Error(t, err)
	assert.Equal(t, 0, totalCount)
	assert.Empty(t, posts)

	// 4. Teardown
}

// 埋め込み用動画URL生成テスト
func TestMakeEmbedMovieURL(t *testing.T) {
	cases := []struct {
		label    string
		movieURL string
		expected string
	}{
		{"非短縮URL、追加パラメータなし", "https://www.youtube.com/watch?v=A1", "https://www.youtube.com/embed/A1"},
		{"非短縮URL、追加パラメータあり", "https://www.youtube.com/watch?v=A1&t=608s", "https://www.youtube.com/embed/A1"},
		{"非短縮URL、動画のキーなし", "https://www.youtube.com/watch?v=", ""},
		{"短縮URL", "https://youtu.be/A1", "https://www.youtube.com/embed/A1"},
		{"短縮URL、動画のキーなし", "https://youtu.be/", ""},
		{"モバイル版", "https://m.youtube.com/watch?v=A1", "https://www.youtube.com/embed/A1"},
		{"HTTP", "http://www.youtube.com/watch?v=A1", "http://www.youtube.com/embed/A1"},
		{"空", "", ""},
		{"形式不正", "dummy", ""},
	}

	for _, test := range cases {
		// 1. Setup

		// 2. Exercise
		embedMovieURL := makeEmbedMovieURL(test.movieURL)

		// 3. Verify
		assert.Equal(t, test.expected, embedMovieURL, test.label)

		// 4. Teardown
	}
}

// 投稿詳細テスト
func TestGetPost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	expected := makeGetPostResult(id)
	repository.On("FetchByID", id).Return(expected, nil)

	// 2. Exercise
	post, err := usecase.GetPost(id)

	// 3. Verify
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, post.ID)
	assert.Equal(t, expected.CreatedAt, post.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, post.UpdatedAt)
	assert.Equal(t, expected.DeletedAt, post.DeletedAt)
	assert.Equal(t, expected.UserID, post.UserID)
	assert.Equal(t, expected.Title, post.Title)
	assert.Equal(t, expected.Speaker, post.Speaker)
	assert.Equal(t, expected.Detail, post.Detail)
	assert.Equal(t, expected.MovieURL, post.MovieURL)
	assert.Equal(t, expected.UserName, post.UserName)

	// 4. Teardown
}

func TestGetPost_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	repository.On("FetchByID", id).Return(nil, errors.New("error"))

	// 2. Execise
	post, err := usecase.GetPost(id)

	// 3. Verify
	assert.Error(t, err)
	assert.Empty(t, post)

	// 4. Teardown
}

// 投稿更新テスト
func TestUpdatePost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := makePostForInput(id)
	repository.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)

	// 2. Exercise
	err := usecase.UpdatePost(id, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestUpdatePost_error(t *testing.T) {
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	post := makePostForInput(id)
	repository.On("Update", mock.AnythingOfType("*model.Post")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.UpdatePost(id, post.Title, post.Speaker, post.Detail, post.MovieURL)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// 投稿削除テスト
func TestDeletePost_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(nil)

	// 2. Exercise
	err := usecase.DeletePost(id)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

func TestDeletePost_error(t *testing.T) {
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	id := 1
	repository.On("Delete", id).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.DeletePost(id)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// お気に入り登録成功
func TestCreateFavorite_success(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	userID := 1
	postID := 1
	favorite := makeFavorite(userID, postID)
	repository.On("CreateFavorite", mock.AnythingOfType("*model.Favorite")).Return(nil)

	// 2. Exercise
	err := usecase.CreateFavorite(favorite.UserID, favorite.PostID)

	// 3. Verify
	assert.NoError(t, err)

	// 4. Teardown
}

// お気に入り登録エラー
func TestCreateFavorite_error(t *testing.T) {
	// 1. Setup
	repository := mockPostRepository{}
	usecase := NewPostUseCase(&repository)
	userID := 1
	postID := 1
	favorite := makeFavorite(userID, postID)
	repository.On("CreateFavorite", mock.AnythingOfType("*model.Favorite")).Return(errors.New("error"))

	// 2. Exercise
	err := usecase.CreateFavorite(favorite.UserID, favorite.PostID)

	// 3. Verify
	assert.Error(t, err)

	// 4. Teardown
}

// TODO お気に入り一覧取得

// TODO お気に入り削除
