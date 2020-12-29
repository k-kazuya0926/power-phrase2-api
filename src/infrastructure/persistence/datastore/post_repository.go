// Package datastore Infra層のリポジトリ
package datastore

import (
	"fmt"

	"github.com/k-kazuya0926/power-phrase2-api/conf"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

// postRepository 構造体
type postRepository struct {
}

// NewPostRepository PostRepositoryを生成する。
func NewPostRepository() repository.PostRepository {
	return &postRepository{}
}

// Create 投稿登録
func (repository *postRepository) Create(post *model.Post) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(post).Error
}

// Fetch 投稿一覧取得。
// キーワード検索を行わない場合はkeywordに空文字を指定する。
// 投稿ユーザーを限定しない場合はpostUserIDに0を指定する。
// ログインユーザーを限定しない場合はloginUserIDに0を指定する。
func (repository *postRepository) Fetch(limit, page int, keyword string, postUserID, loginUserID int) (totalCount int, posts []*model.GetPostResult, err error) {
	db := conf.NewDBConnection()
	defer db.Close()

	countDb := conf.NewDBConnection()
	defer countDb.Close()

	offset := limit * (page - 1)

	if keyword != "" {
		// キーワードがタイトル、発言者、詳細のいずれかに含まれる
		countDb = countDb.Where("title LIKE ?", "%"+keyword+"%").Or("speaker LIKE ?", "%"+keyword+"%").Or("detail LIKE ?", "%"+keyword+"%")
		db = db.Where("title LIKE ?", "%"+keyword+"%").Or("speaker LIKE ?", "%"+keyword+"%").Or("detail LIKE ?", "%"+keyword+"%")
	}

	if postUserID > 0 { // ユーザーIDが指定されている場合
		countDb = countDb.Where("posts.user_id = ?", postUserID)
		db = db.Where("posts.user_id = ?", postUserID)
	}

	// 投稿総件数取得
	if err = countDb.Model(&model.Post{}).Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	// 投稿一覧取得
	if err = db.Table("posts").
		Select(`posts.*,
			users.name as user_name,
			users.image_file_path as user_image_file_path,
			(SELECT count(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL) AS comment_count,
			(CASE WHEN favorites.id IS NULL THEN false ELSE true END) AS is_favorite,
			(SELECT count(*) FROM favorites AS f WHERE f.post_id = posts.id) AS favorite_count
		`).
		Joins(fmt.Sprintf(`JOIN users ON users.id = posts.user_id AND users.deleted_at IS NULL
			LEFT JOIN favorites ON favorites.post_id = posts.id AND favorites.user_id = %d`, loginUserID)).
		Order("posts.id DESC").Limit(limit).Offset(offset).
		Find(&posts).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, posts, err
}

// FetchByID 投稿1件取得
func (repository *postRepository) FetchByID(id, loginUserID int) (*model.GetPostResult, error) {
	db := conf.NewDBConnection()
	defer db.Close()

	post := model.GetPostResult{}
	post.ID = id
	if err := db.Table("posts").
		Select(`posts.*,
			users.name as user_name,
			users.image_file_path as user_image_file_path,
			(SELECT count(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL) AS comment_count,
			(CASE WHEN favorites.id IS NULL THEN false ELSE true END) AS is_favorite,
			(SELECT count(*) FROM favorites AS f WHERE f.post_id = posts.id) AS favorite_count
		`).
		Joins(fmt.Sprintf(`JOIN users ON users.id = posts.user_id AND users.deleted_at IS NULL
			LEFT JOIN favorites ON favorites.post_id = posts.id AND favorites.user_id = %d`, loginUserID)).
		First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

// Update 投稿更新
func (repository *postRepository) Update(u *model.Post) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Model(u).Update(u).Error
}

// Delete 投稿削除
func (repository *postRepository) Delete(id int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	post := model.Post{ID: id}
	return db.Delete(&post).Error
}

// CreateComment コメント登録
func (repository *postRepository) CreateComment(comment *model.Comment) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(comment).Error
}

// FetchComments コメント一覧取得
func (repository *postRepository) FetchComments(postID, limit, page int) (totalCount int, comments []*model.GetCommentResult, err error) {
	db := conf.NewDBConnection()
	defer db.Close()

	countDb := conf.NewDBConnection()
	defer countDb.Close()

	if err = countDb.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	offset := limit * (page - 1)
	if err = db.Table("comments").
		Select("comments.*, users.name as user_name, users.image_file_path as user_image_file_path").
		Joins("JOIN users on users.id = comments.user_id AND users.deleted_at IS NULL").
		Where("post_id = ?", postID).
		Order("id DESC").Limit(limit).Offset(offset).
		Find(&comments).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, comments, err
}

// DeleteComment コメント削除
func (repository *postRepository) DeleteComment(id int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	comment := model.Comment{ID: id}
	return db.Delete(&comment).Error
}

// CreateFavorite お気に入り登録
func (repository *postRepository) CreateFavorite(favorite *model.Favorite) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Create(favorite).Error
}

// FetchFavorites お気に入り一覧取得
func (repository *postRepository) FetchFavorites(userID, limit, page int) (totalCount int, posts []*model.GetPostResult, err error) {
	db := conf.NewDBConnection()
	defer db.Close()

	countDb := conf.NewDBConnection()
	defer countDb.Close()

	if err = countDb.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return 0, nil, err
	}

	offset := limit * (page - 1)
	if err = db.Unscoped().Table("favorites").
		Select(`posts.*,
			users.name AS user_name,
			users.image_file_path AS user_image_file_path,
			(SELECT count(*) FROM comments WHERE comments.post_id = posts.id AND comments.deleted_at IS NULL) AS comment_count,
			true AS is_favorite
		`).
		Joins(`JOIN posts ON posts.id = favorites.post_id AND posts.deleted_at IS NULL
			JOIN users ON users.id = posts.user_id AND users.deleted_at IS NULL`).
		Where("favorites.user_id = ?", userID).
		Order("posts.id DESC").Limit(limit).Offset(offset).
		Find(&posts).Error; err != nil {
		return 0, nil, err
	}

	return totalCount, posts, err
}

// DeleteFavorite お気に入り削除
func (repository *postRepository) DeleteFavorite(userID, postID int) error {
	db := conf.NewDBConnection()
	defer db.Close()

	return db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Favorite{}).Error
}
