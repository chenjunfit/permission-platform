package jwt

import (
	"fmt"
	"testing"
)

func TestGenerater(t *testing.T) {
	token := New("permission_platform_key", "permission-platform")
	claim := MapClaims{
		"biz_id": 123456,
	}
	tokenStr, _ := token.Encode(claim)
	fmt.Println(tokenStr)
}
