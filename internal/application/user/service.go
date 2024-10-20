package user

import (
	"blog-api/internal/domain/user"
	"blog-api/internal/infrastructure/auth"
	"blog-api/internal/infrastructure/hash"
	"errors"
	"time"
)

// Service 封裝了用戶相關的業務邏輯
type Service struct {
	repo       user.Repository
	jwtService *auth.JWTService
}

// NewService 創建一個新的用戶服務實例
func NewService(repo user.Repository, jwtService *auth.JWTService) *Service {
	return &Service{repo: repo, jwtService: jwtService}
}

// RegisterInput 定義註冊所需的輸入數據
type RegisterInput struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Register 處理用戶註冊邏輯
func (s *Service) Register(input RegisterInput) error {
	// 檢查用戶名是否已存在
	if _, err := s.repo.FindByUsername(input.Username); err == nil {
		return user.ErrDuplicateUsername
	}

	// 檢查郵箱是否已存在
	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return user.ErrDuplicateEmail
	}

	// 驗證密碼
	if err := user.ValidatePassword(input.Password); err != nil {
		return err
	}

	// 對密碼進行哈希處理
	hashedPassword, err := hash.GenerateFromPassword(input.Password, hash.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &user.User{
		Username:          input.Username,
		Email:             input.Email,
		PasswordHash:      string(hashedPassword),
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		IsActive:          true,
		PasswordChangedAt: time.Now(), // 設置初始密碼修改時間
	}

	return s.repo.Create(newUser)
}

// LoginInput 定義登錄所需的輸入數據
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 處理用戶登錄邏輯
func (s *Service) Login(input LoginInput) (string, error) {
	u, err := s.repo.FindByUsername(input.Username)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", user.ErrInvalidPassword // 為了安全，不透露用戶不存在的信息
		}
		return "", err
	}

	if !u.IsActive {
		return "", errors.New("account is not active")
	}

	if err := hash.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)); err != nil {
		return "", user.ErrInvalidPassword
	}

	// 更新最後登錄時間
	u.UpdateLastLogin()
	if err := s.repo.Update(u); err != nil {
		return "", err
	}

	// 生成 JWT 令牌，包含密碼修改時間
	token, err := s.jwtService.GenerateToken(u.ID, u.PasswordChangedAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserProfile 根據用戶ID獲取用戶信息
func (s *Service) GetUserProfile(id uint) (*user.User, error) {
	return s.repo.FindByID(id)
}

// ChangePasswordInput 定義更改密碼所需的輸入數據
type ChangePasswordInput struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

// ChangePassword 處理更改密碼的邏輯
func (s *Service) ChangePassword(userID uint, input ChangePasswordInput) error {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// 驗證當前密碼
	if err := hash.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// 驗證新密碼
	if err := user.ValidatePassword(input.NewPassword); err != nil {
		return err
	}

	// 對新密碼進行哈希處理
	hashedPassword, err := hash.GenerateFromPassword(input.NewPassword, hash.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密碼和密碼修改時間
	u.PasswordHash = string(hashedPassword)
	u.PasswordChangedAt = time.Now()
	return s.repo.Update(u)
}
