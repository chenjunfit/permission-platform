package ioc

import "github.com/permission-dev/internal/pkg/jwt"

func InitJWTToken() *jwt.Token {
	return jwt.New("permission_platform_key", "permission-platform")
}
