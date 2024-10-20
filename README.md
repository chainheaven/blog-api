# Blog API
## 介紹
這是一個基於 Go 語言開發的 Blog 後端 API，遵循領域驅動設計（DDD）架構。

### 功能特點

- 用戶註冊和登錄
- JWT 認證
- 文章的創建、讀取、更新和刪除（CRUD）操作
- 密碼加密存儲
- 分頁獲取文章列表

## 技術棧

- Go
- Gin Web 框架
- GORM ORM 庫
- PostgreSQL 數據庫
- JWT 認證
- Swagger 用於 API 文檔

## 安裝

1. 克隆專案：
```git clone https://github.com/your-username/blog-api.git```

2. 進入專案目錄：
```cd blog-api```

3. 安裝依賴：
```go mod download```
```go mod tidy```


## 環境設置

1. 複製 .env.example 文件並重命名為 .env：
Copycp .env.example .env

2. 打開 .env 文件並填入您的實際配置值：

將 your_username、your_password 和 your_database_name 替換為您的 PostgreSQL 數據庫憑證。
將 your_jwt_secret_key 替換為一個安全的隨機字符串。
如果需要，可以修改 PORT 值。

## 生成 API 文檔
運行以下命令生成 Swagger 文檔：
```swag init -g cmd/api/main.go```
## 運行應用
使用以下命令啟動應用：
```go run cmd/api/main.go```
## API 文檔
啟動應用後，Swagger UI 可以在 http://localhost:8080/swagger/index.html 訪問。
