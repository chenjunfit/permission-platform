package ioc

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/permission-dev/internal/pkg/jwt"
)

func InitJwtToken() *jwt.Token {
	type Config struct {
		Key    string `json:"key"`
		Issuer string `json:"issuer"`
	}
	var cfg Config
	err := econf.UnmarshalKey("jwt", &cfg)
	if err != nil {
		panic(err)
	}
	return jwt.New(cfg.Key, cfg.Issuer)
}
