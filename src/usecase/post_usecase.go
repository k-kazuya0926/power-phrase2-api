package usecase

import (
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
	"github.com/k-kazuya0926/power-phrase2-api/domain/repository"
)

type PostUseCase interface {
	CreatePost(userID int, title, speaker, detail, movieURL string) (err error)
	GetPosts(limit, offset int, keyword string) (totalCount int, posts []*model.GetPostResult, err error)
	GetPost(id int) (*model.Post, error)
	UpdatePost(ID int, title, speaker, detail, movieURL string) error
	DeletePost(id int) error
}

type postUseCase struct {
	repository.PostRepository
}

// NewPostUseCase PostUseCaseを取得します.
func NewPostUseCase(repository repository.PostRepository) PostUseCase {
	return &postUseCase{repository}
}

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

func (usecase *postUseCase) GetPosts(limit, page int, keyword string) (totalCount int, posts []*model.GetPostResult, err error) {
	totalCount, posts, err = usecase.PostRepository.Fetch(limit, page, keyword)
	if err != nil {
		return 0, nil, err
	}
	return totalCount, posts, nil
}

func (usecase *postUseCase) GetPost(id int) (*model.Post, error) {
	post, err := usecase.PostRepository.FetchByID(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

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

func (usecase *postUseCase) DeletePost(id int) error {
	if err := usecase.PostRepository.Delete(id); err != nil {
		return err
	}
	return nil
}
