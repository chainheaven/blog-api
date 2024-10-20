package post

import (
	"errors"
	"time"
)

// Post 文章
type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string    `json:"title" binding:"required" gorm:"type:varchar(255);not null"`
	Content   string    `json:"content" binding:"required" gorm:"type:text;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

// 定義一些常見的錯誤
var (
	ErrPostNotFound   = errors.New("post not found")
	ErrUnauthorized   = errors.New("unauthorized to modify this post")
	ErrInvalidTitle   = errors.New("invalid post title")
	ErrInvalidContent = errors.New("invalid post content")
)

// Repository 定義文章存儲的接口
type Repository interface {
	FindAll(page, pageSize int) ([]Post, error)
	FindByID(id uint) (*Post, error)
	Create(post *Post) error
	Update(post *Post) error
	Delete(id uint) error
}

// ValidateTitle 驗證文章標題是否符合要求
func ValidateTitle(title string) error {
	if len(title) < 3 || len(title) > 255 {
		return ErrInvalidTitle
	}
	return nil
}

// ValidateContent 驗證文章內容是否符合要求
func ValidateContent(content string) error {
	if len(content) < 10 {
		return ErrInvalidContent
	}
	return nil
}

// IsAuthor 檢查給定的用戶ID是否為文章作者
func (p *Post) IsAuthor(userID uint) bool {
	return p.UserID == userID
}

// UpdateContent 更新文章內容
func (p *Post) UpdateContent(newTitle, newContent string) error {
	if err := ValidateTitle(newTitle); err != nil {
		return err
	}
	if err := ValidateContent(newContent); err != nil {
		return err
	}
	p.Title = newTitle
	p.Content = newContent
	p.UpdatedAt = time.Now()
	return nil
}
