// Package usecase Application Service層。
package usecase

import (
	"net/url"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// PostUseCase インターフェース
type PostUseCase interface {
	CreatePost(userID int, title, speaker, detail, movieURL string) (err error)
	GetPosts(limit, offset int, keyword string, userID int) (totalCount int, posts []*model.GetPostResult, err error)
	GetPost(id int) (*model.GetPostResult, error)
	UpdatePost(ID int, title, speaker, detail, movieURL string) error
	DeletePost(id int) error
}

// postUseCase 構造体
type postUseCase struct {
	repository.PostRepository
}

// NewPostUseCase PostUseCaseを生成。
func NewPostUseCase(repository repository.PostRepository) PostUseCase {
	return &postUseCase{repository}
}

// CreatePost 登録
func (usecase *postUseCase) CreatePost(userID int, title, speaker, detail, movieURL string) (err error) {
	post := model.Post{
		UserID:   userID,
		Title:    title,
		Speaker:  speaker,
		Detail:   detail,
		MovieURL: movieURL,
	}
	err = usecase.PostRepository.Create(&post)

	return err
}

// GetPosts 一覧取得。キーワード検索を行わない場合はkeywordに空文字を指定する。ユーザーを限定しない場合はuserIDに0を指定する。
func (usecase *postUseCase) GetPosts(limit, page int, keyword string, userID int) (totalCount int, posts []*model.GetPostResult, err error) {
	totalCount, posts, err = usecase.PostRepository.Fetch(limit, page, keyword, userID)
	if err != nil {
		return 0, nil, err
	}

	// 動画URL加工
	for _, post := range posts {
		post.EmbedMovieURL = makeEmbedMovieURL(post.MovieURL)
	}

	return totalCount, posts, nil
}

// makeEmbedMovieURL 埋め込み用動画URLを生成する。
func makeEmbedMovieURL(movieURL string) string {
	u, err := url.Parse(movieURL)
	if err != nil {
		return ""
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return ""
	}
	v, ok := m["v"]
	if ok {
		return u.Scheme + "://" + u.Host + "/embed/" + v[0]
	} else {
		return ""
	}
}

// GetPost 1件取得
func (usecase *postUseCase) GetPost(id int) (*model.GetPostResult, error) {
	post, err := usecase.PostRepository.FetchByID(id)
	if err != nil {
		return nil, err
	}

	// 動画URL加工
	post.EmbedMovieURL = makeEmbedMovieURL(post.MovieURL)

	return post, nil
}

// UpdatePost 更新
func (usecase *postUseCase) UpdatePost(ID int, title, speaker, detail, movieURL string) error {
	post := model.Post{
		ID:       ID,
		Title:    title,
		Speaker:  speaker,
		Detail:   detail,
		MovieURL: movieURL,
	}
	if err := usecase.PostRepository.Update(&post); err != nil {
		return err
	}
	return nil
}

// DeletePost 削除
func (usecase *postUseCase) DeletePost(id int) error {
	if err := usecase.PostRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
