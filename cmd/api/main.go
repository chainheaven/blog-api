package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "blog-api/docs"
	"blog-api/internal/application/post"
	"blog-api/internal/application/user"
	"blog-api/internal/infrastructure/auth"
	"blog-api/internal/infrastructure/http"
	"blog-api/internal/infrastructure/http/handlers"
	"blog-api/internal/infrastructure/postgres"

	"github.com/joho/godotenv"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Blog API
// @version 1.0
// @description This is a sample server for a blog API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 加載 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 從環境變量獲取數據庫連接信息
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	// 配置數據庫連接
	db, err := gorm.Open(pgDriver.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化存儲層
	userRepo := postgres.NewUserRepository(db)
	postRepo := postgres.NewPostRepository(db)

	// 初始化 JWT 服務
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET_KEY"))

	// 初始化服務層
	userService := user.NewService(userRepo, jwtService)
	postService := post.NewService(postRepo)

	// 初始化處理器
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	// 設置路由
	r := http.SetupRouter(userHandler, postHandler, jwtService, userService)

	// 獲取服務器端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Swagger UI URL
	fmt.Printf("\nSwagger UI is available at: http://localhost:%s/swagger/index.html\n\n", port)

	// 啟動服務器
	log.Printf("Server is starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
