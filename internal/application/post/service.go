package post

import (
	"blog-api/internal/domain/post"
)

// Service 封裝了文章相關的業務邏輯
type Service struct {
	repo post.Repository
}

// NewService 創建一個新的文章服務實例
func NewService(repo post.Repository) *Service {
	return &Service{repo: repo}
}

// GetPosts 獲取文章列表
func (s *Service) GetPosts(page int) ([]post.Post, error) {
	return s.repo.FindAll(page, 10) // 10 posts per page
}

// GetPostByID 根據ID獲取單個文章
func (s *Service) GetPostByID(id uint) (*post.Post, error) {
	return s.repo.FindByID(id)
}

// CreatePost 創建新文章
func (s *Service) CreatePost(p *post.Post) error {
	if err := post.ValidateTitle(p.Title); err != nil {
		return err
	}
	if err := post.ValidateContent(p.Content); err != nil {
		return err
	}
	return s.repo.Create(p)
}

// UpdatePost 更新現有文章
func (s *Service) UpdatePost(p *post.Post, userID uint) error {
	existingPost, err := s.repo.FindByID(p.ID)
	if err != nil {
		return err
	}
	if !existingPost.IsAuthor(userID) {
		return post.ErrUnauthorized
	}
	if err := existingPost.UpdateContent(p.Title, p.Content); err != nil {
		return err
	}
	return s.repo.Update(existingPost)
}

// DeletePost 刪除文章
func (s *Service) DeletePost(id, userID uint) error {
	existingPost, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !existingPost.IsAuthor(userID) {
		return post.ErrUnauthorized
	}
	return s.repo.Delete(id)
}
