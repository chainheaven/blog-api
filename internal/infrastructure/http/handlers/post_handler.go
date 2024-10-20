package handlers

import (
	appPost "blog-api/internal/application/post"
	"blog-api/internal/domain/post"
	"blog-api/internal/infrastructure/http/middlewares"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PostInput 用於接收用戶輸入的文章數據
// @Description 用於創建或更新文章的輸入模型
type PostInput struct {
	Title   string `json:"title" binding:"required" example:"My Blog Post"`
	Content string `json:"content" binding:"required" example:"This is the content of my blog post."`
}

type PostHandler struct {
	postService *appPost.Service
}

func NewPostHandler(postService *appPost.Service) *PostHandler {
	return &PostHandler{postService: postService}
}

// GetPosts 返回文章列表
// @Summary 獲取文章列表
// @Description 返回分頁的文章列表，每頁10篇
// @Tags posts
// @Produce json
// @Param page query int false "頁碼" default(1)
// @Success 200 {array} post.Post
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	posts, err := h.postService.GetPosts(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPost 返回單篇文章
// @Summary 獲取文章詳情
// @Description 根據ID返回單篇文章的詳細內容
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} post.Post
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	post, err := h.postService.GetPostByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// CreatePost 創建新文章
// @Summary 創建新文章
// @Description 創建一篇新文章，需要用戶登錄
// @Tags posts
// @Accept json
// @Produce json
// @Param post body PostInput true "文章內容"
// @Security BearerAuth
// @Success 201 {object} post.Post
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middlewares.GetUserID(c)
	now := time.Now()
	newPost := &post.Post{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.postService.CreatePost(newPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

// UpdatePost 更新文章
// @Summary 更新文章
// @Description 更新現有文章，需要用戶登錄且為作者
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param post body PostInput true "更新的文章內容"
// @Security BearerAuth
// @Success 200 {object} post.Post
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input PostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middlewares.GetUserID(c)

	updatedPost := &post.Post{
		ID:        uint(id),
		Title:     input.Title,
		Content:   input.Content,
		UserID:    userID,
		UpdatedAt: time.Now(),
	}

	if err := h.postService.UpdatePost(updatedPost, userID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

// DeletePost 刪除文章
// @Summary 刪除文章
// @Description 刪除指定文章，需要用戶登錄且為作者
// @Tags posts
// @Produce json
// @Param id path int true "文章ID"
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID, _ := middlewares.GetUserID(c)

	if err := h.postService.DeletePost(uint(id), userID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
