package http

import (
	"blog-api/internal/application/user"
	"blog-api/internal/infrastructure/auth"
	"blog-api/internal/infrastructure/http/handlers"
	"blog-api/internal/infrastructure/http/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 配置 API 路由
func SetupRouter(userHandler *handlers.UserHandler, postHandler *handlers.PostHandler, jwtService *auth.JWTService, userService *user.Service) *gin.Engine {
	r := gin.Default()

	// API 路由
	api := r.Group("/api/v1")
	{
		// 用戶相關路由
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// 文章相關路由
		posts := api.Group("/posts")
		{
			posts.GET("", postHandler.GetPosts)
			posts.GET("/:id", postHandler.GetPost)

			// 需要認證的路由
			authorized := posts.Group("/")
			authorized.Use(middlewares.AuthMiddleware(jwtService, userService))
			{
				authorized.POST("", postHandler.CreatePost)
				authorized.PUT("/:id", postHandler.UpdatePost)
				authorized.DELETE("/:id", postHandler.DeletePost)
			}
		}

		// 用戶認證路由
		authorized := api.Group("/")
		authorized.Use(middlewares.AuthMiddleware(jwtService, userService))
		{
			authorized.GET("/profile", userHandler.GetProfile)
			authorized.POST("/change-password", userHandler.ChangePassword)
		}
	}

	// Swagger 文檔路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
