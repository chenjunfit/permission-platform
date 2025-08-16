package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/permission-dev/internal/errs"
	"strings"
	"time"
)

// 工具包的错误信息暴露
var (
	ErrTokenExpired          = jwt.ErrTokenExpired
	ErrTokenSignatureInvaild = jwt.ErrSignatureInvalid
)

// type 定义不会继承原来类型的方法
// type =会继承
type MapClaims jwt.MapClaims

type Token struct {
	key    string //key 签发的key
	issuer string //签发者
}

func New(key, issuer string) *Token {
	return &Token{
		key:    key,
		issuer: issuer,
	}
}
func (t *Token) Decode(tokenString string) (MapClaims, error) {
	//去除可能Bearer前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w,%v", errs.ErrSupportedSignAlgorithm, token.Header["alg"])
		}
		return []byte(t.key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w %v", errs.ErrDecodeJWTTokenFailed, err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return MapClaims(claims), nil
	}
	return nil, fmt.Errorf("%w", errs.ErrInvalidJWTToken)
}
func (t *Token) Encode(clustomClaims MapClaims) (string, error) {
	mapClaims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"iss": t.issuer,
	}

	//合并用户自定义声明，含默认的
	for k, v := range clustomClaims {
		mapClaims[k] = v
	}
	//自动处理过期时间
	const day = 24 * time.Hour
	if _, ok := mapClaims["exp"]; !ok {
		mapClaims["exp"] = time.Now().Add(day).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(t.key))
}
