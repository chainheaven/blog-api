package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// 定義常量
const (
	TokenExpireDuration = time.Hour * 24 // 令牌有效期為 24 小時
)

// Claims 自定義 JWT 聲明結構體
type Claims struct {
	UserID            uint      `json:"user_id"`
	PasswordChangedAt time.Time `json:"pwd_changed_at"`
	jwt.StandardClaims
}

// 定義錯誤
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrNoSecretKey  = errors.New("no secret key set")
)

// JWTService 提供 JWT 相關功能
type JWTService struct {
	secretKey []byte
}

// NewJWTService 創建一個新的 JWTService 實例
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{secretKey: []byte(secretKey)}
}

// GenerateToken 生成 JWT 令牌
func (s *JWTService) GenerateToken(userID uint, passwordChangedAt time.Time) (string, error) {
	claims := Claims{
		UserID:            userID,
		PasswordChangedAt: passwordChangedAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "blog-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ParseToken 解析 JWT 令牌
func (s *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrInvalidToken
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateToken 驗證 JWT 令牌
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	return s.ParseToken(tokenString)
}
