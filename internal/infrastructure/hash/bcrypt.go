package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 定義一些常量
const (
	MinCost     int = 4  // bcrypt 的最小代價
	MaxCost     int = 31 // bcrypt 的最大代價
	DefaultCost int = 10 // 默認的代價
)

// 定義可能的錯誤
var (
	ErrInvalidCost = errors.New("invalid bcrypt cost")
	ErrHashFailed  = errors.New("hashing password failed")
)

// GenerateFromPassword 使用 bcrypt 對密碼進行哈希
func GenerateFromPassword(password string, cost int) ([]byte, error) {
	// 驗證 cost 值
	if cost < MinCost || cost > MaxCost {
		return nil, ErrInvalidCost
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, ErrHashFailed
	}

	return hash, nil
}

// CompareHashAndPassword 比較哈希值和密碼
func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// Cost 返回給定哈希的代價
func Cost(hashedPassword []byte) (int, error) {
	return bcrypt.Cost(hashedPassword)
}

// GenerateFromPasswordWithDefault 使用默認代價生成密碼哈希
func GenerateFromPasswordWithDefault(password string) ([]byte, error) {
	return GenerateFromPassword(password, DefaultCost)
}

// IsHashedPassword 檢查給定的字節序列是否可能是 bcrypt 哈希
func IsHashedPassword(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return data[0] == '$' && data[1] == '2' && (data[2] == 'a' || data[2] == 'b' || data[2] == 'y') && data[3] == '$'
}
