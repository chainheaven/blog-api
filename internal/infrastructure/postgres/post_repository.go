package postgres

import (
	"blog-api/internal/domain/post"

	"gorm.io/gorm"
)

// PostRepository 實現 post.Repository 接口
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 創建一個新的 PostRepository 實例
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// FindAll 獲取分頁的文章列表
func (r *PostRepository) FindAll(page, pageSize int) ([]post.Post, error) {
	var posts []post.Post
	offset := (page - 1) * pageSize
	err := r.db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, err
}

// FindByID 根據ID查找文章
func (r *PostRepository) FindByID(id uint) (*post.Post, error) {
	var p post.Post
	if err := r.db.First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, post.ErrPostNotFound
		}
		return nil, err
	}
	return &p, nil
}

// Create 創建新文章
func (r *PostRepository) Create(post *post.Post) error {
	return r.db.Create(post).Error
}

// Update 更新現有文章
func (r *PostRepository) Update(post *post.Post) error {
	return r.db.Save(post).Error
}

// Delete 刪除文章
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&post.Post{}, id).Error
}
