package postgres

import (
	"blog-api/internal/domain/user"

	"gorm.io/gorm"
)

// UserRepository 實現 user.Repository 接口，使用 GORM 和 PostgreSQL
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 創建一個新的 UserRepository 實例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 將新用戶保存到數據庫
func (r *UserRepository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

// FindByUsername 根據用戶名從數據庫中查找用戶
func (r *UserRepository) FindByUsername(username string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByID 根據ID從數據庫中查找用戶
func (r *UserRepository) FindByID(id uint) (*user.User, error) {
	var u user.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByEmail 根據郵箱從數據庫中查找用戶
func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Update 更新數據庫中的用戶信息
func (r *UserRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}

// Delete 從數據庫中刪除用戶
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&user.User{}, id).Error
}
