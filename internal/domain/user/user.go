package user

import (
	"errors"
	"time"
)

// User 代表系統中的用戶實體
type User struct {
	ID                uint       `json:"id" gorm:"primaryKey" example:"1"`
	Username          string     `json:"username" gorm:"uniqueIndex;not null" example:"johndoe"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null" example:"john@example.com"`
	PasswordHash      string     `json:"-" gorm:"not null"` // 使用 "-" 標籤避免在 JSON 響應中返回密碼哈希
	PasswordChangedAt time.Time  `json:"passwordChangedAt" gorm:"not null" example:"2024-10-20T15:00:00Z"`
	FirstName         string     `json:"firstName" example:"John"`
	LastName          string     `json:"lastName" example:"Doe"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP" example:"2024-10-20T14:00:00Z"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP" example:"2024-10-20T14:30:00Z"`
	LastLogin         *time.Time `json:"lastLogin,omitempty" example:"2024-10-20T16:00:00Z"`
	IsActive          bool       `json:"isActive" gorm:"default:true" example:"true"`
}

// 定義一些常見的錯誤
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
)

// Repository 定義了用戶資料持久化的接口
type Repository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// ValidatePassword 驗證密碼是否符合要求
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

// FullName 返回用戶的全名
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// UpdateLastLogin 更新用戶的最後登錄時間
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
}

// ChangePassword 更改用戶密碼
func (u *User) ChangePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.UpdatedAt = time.Now()
}
