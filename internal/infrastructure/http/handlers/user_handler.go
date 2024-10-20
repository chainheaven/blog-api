package handlers

import (
	"blog-api/internal/application/user"
	"blog-api/internal/infrastructure/http/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 處理與用戶相關的 HTTP 請求
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler 創建一個新的 UserHandler 實例
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 處理用戶註冊請求
// @Summary 註冊新用戶
// @Description 創建一個新的用戶帳戶
// @Tags user
// @Accept  json
// @Produce  json
// @Param   input body user.RegisterInput true "註冊信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var input user.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Register(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login 處理用戶登錄請求
// @Summary 用戶登錄
// @Description 驗證用戶憑證並返回 JWT 令牌
// @Tags user
// @Accept  json
// @Produce  json
// @Param   input body user.LoginInput true "登錄信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var input user.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetProfile 獲取用戶資料
// @Summary 獲取用戶資料
// @Description 獲取當前登錄用戶的資料
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := middlewares.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword 更改用戶密碼
// @Summary 更改用戶密碼
// @Description 更改當前登錄用戶的密碼
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body user.ChangePasswordInput true "更改密碼信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var input user.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := middlewares.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.userService.ChangePassword(userID, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
