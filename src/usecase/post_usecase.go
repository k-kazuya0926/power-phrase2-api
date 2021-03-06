// Package usecase Application Service層。
package usecase

import (
	"net/url"

	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// PostUseCase インターフェース
type PostUseCase interface {
	// 投稿登録
	CreatePost(userID int, title, speaker, detail, movieURL string) (err error)
	// 投稿一覧取得
	GetPosts(limit, offset int, keyword string, postUserID, loginUserID int) (totalCount int, posts []*model.GetPostResult, err error)
	// 投稿詳細取得
	GetPost(id, loginUserID int) (*model.GetPostResult, error)
	// 投稿更新
	UpdatePost(ID int, title, speaker, detail, movieURL string) error
	// 投稿削除
	DeletePost(id int) error

	// お気に入り登録
	CreateFavorite(postID, userID int) (err error)
	// お気に入り一覧取得
	GetFavorites(userID, limit, offset int) (totalCount int, posts []*model.GetPostResult, err error)
	// お気に入り削除
	DeleteFavorite(userID, postID int) error
}

// postUseCase 構造体
type postUseCase struct {
	repository.PostRepository
}

// NewPostUseCase PostUseCaseを生成。
func NewPostUseCase(repository repository.PostRepository) PostUseCase {
	return &postUseCase{repository}
}

// CreatePost 投稿登録
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

// GetPosts 一覧取得。
// キーワード検索を行わない場合はkeywordに空文字を指定する。
// 投稿ユーザーを限定しない場合はpostUserIDに0を指定する。
// ログインユーザーを限定しない場合はloginUserIDに0を指定する。
func (usecase *postUseCase) GetPosts(limit, page int, keyword string, postUserID, loginUserID int) (totalCount int, posts []*model.GetPostResult, err error) {
	totalCount, posts, err = usecase.PostRepository.Fetch(limit, page, keyword, postUserID, loginUserID)
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
	urlStruct, err := url.Parse(movieURL)
	if err != nil {
		return ""
	}

	// 短縮URL(例：https://youtu.be/9LmL92WLgfc)の場合
	if urlStruct.Hostname() == "youtu.be" {
		// 動画のキーが含まれていない場合
		if urlStruct.Path == "/" {
			return ""
		}

		return urlStruct.Scheme + "://www.youtube.com/embed" + urlStruct.Path
	}

	// 以下、短縮URLでない場合

	parameterMap, err := url.ParseQuery(urlStruct.RawQuery)
	if err != nil {
		return ""
	}
	// 動画のキー
	v, ok := parameterMap["v"]
	if ok && v[0] != "" {
		// ホスト名が「m.youtube.com」などである場合も「www.youtube.com」にする
		return urlStruct.Scheme + "://www.youtube.com/embed/" + v[0]
	}

	return ""
}

// GetPost 投稿詳細取得
func (usecase *postUseCase) GetPost(id, loginUserID int) (*model.GetPostResult, error) {
	post, err := usecase.PostRepository.FetchByID(id, loginUserID)
	if err != nil {
		return nil, err
	}

	// 動画URL加工
	post.EmbedMovieURL = makeEmbedMovieURL(post.MovieURL)

	return post, nil
}

// UpdatePost 投稿更新
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

// DeletePost 投稿削除
func (usecase *postUseCase) DeletePost(id int) error {
	if err := usecase.PostRepository.Delete(id); err != nil {
		return err
	}
	return nil
}

// CreateFavorite お気に入り登録
func (usecase *postUseCase) CreateFavorite(userID, postID int) (err error) {
	favorite := model.Favorite{
		UserID: userID,
		PostID: postID,
	}
	err = usecase.PostRepository.CreateFavorite(&favorite)

	return err
}

// GetFavorites お気に入り一覧取得
func (usecase *postUseCase) GetFavorites(userID, limit, page int) (totalCount int, posts []*model.GetPostResult, err error) {
	totalCount, posts, err = usecase.PostRepository.FetchFavorites(userID, limit, page)
	if err != nil {
		return 0, nil, err
	}

	// 動画URL加工
	for _, post := range posts {
		post.EmbedMovieURL = makeEmbedMovieURL(post.MovieURL)
	}

	return totalCount, posts, nil
}

// DeleteFavorite お気に入り削除
func (usecase *postUseCase) DeleteFavorite(userID, postID int) error {
	if err := usecase.PostRepository.DeleteFavorite(userID, postID); err != nil {
		return err
	}
	return nil
}
