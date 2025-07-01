package jwt

import (
	"errors"
	"time"

	"todolist/config"

	"github.com/golang-jwt/jwt"
)

// 自定义错误
var (
	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	ErrTokenMalformed   = errors.New("令牌格式错误")
	ErrTokenInvalid     = errors.New("令牌无效")
)

// CustomClaims 自定义 JWT 声明
type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID int, username string) (string, error) {
	// 获取配置
	jwtConfig := config.GlobalConfig.JWT

	// 创建 claims
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtConfig.ExpireHours).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    jwtConfig.Issuer,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获得完整的编码后的字符串令牌
	return token.SignedString([]byte(jwtConfig.SecretKey))
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.SecretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, ErrTokenMalformed
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, ErrTokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, ErrTokenNotValidYet
			default:
				return nil, ErrTokenInvalid
			}
		}
		return nil, err
	}

	// 验证令牌
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ValidateToken 验证令牌是否有效
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

// GetUserIDFromToken 从令牌中获取用户ID
func GetUserIDFromToken(tokenString string) (int, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetUsernameFromToken 从令牌中获取用户名
func GetUsernameFromToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Username, nil
}
