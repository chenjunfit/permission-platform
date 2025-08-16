package jwt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	defaultIssuer = "permission-dev"
)

func TestJwtEncode(t *testing.T) {
	key := "permission-dev"
	token := New(key, defaultIssuer)
	//测试场景
	testcases := []struct {
		name        string
		customClaim MapClaims
		wantErr     bool
	}{
		{
			name: "测试带用户id的令牌",
			customClaim: MapClaims{
				"biz_id":  float64(1),
				"user_id": float64(999),
			},
			wantErr: false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := token.Encode(tt.customClaim)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			claims, err := token.Decode(tokenString)
			assert.NoError(t, err)
			assert.NotEmpty(t, claims["exp"])
			assert.NotEmpty(t, claims["iat"])
			assert.Equal(t, defaultIssuer, claims["iss"])
			//验证自定义声明
			for k, v := range tt.customClaim {
				fmt.Println(v)
				assert.Equal(t, v, claims[k])
			}
		})
	}
}

func TestJwtDecode(t *testing.T) {
	key := "permission-dev2"
	token1 := New(key, defaultIssuer)

	validClaim := MapClaims{
		"user_id": "123456",
		"role":    "admin",
	}
	validToken, err := token1.Encode(validClaim)
	assert.NoError(t, err)

	expiredClaim := MapClaims{
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}
	expiredToken, err := token1.Encode(expiredClaim)
	assert.NoError(t, err)

	//测试集合
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有效的令牌解析",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "过期令牌解析",
			token:   expiredToken,
			wantErr: true,
		},
		{
			name:    "带Bearer前缀的有效令牌",
			token:   "Bearer " + validToken,
			wantErr: false,
		},
		{
			name:    "无效令牌格式",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "空令牌",
			token:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		claim, err := token1.Decode(tt.token)
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, claim)

			if !tt.wantErr && tt.token == validToken || tt.token == "Bearer "+validToken {
				assert.Equal(t, "123456", claim["user_id"])
				assert.Equal(t, "admin", claim["role"])
			}
		})
	}
}
