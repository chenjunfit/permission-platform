package auth

import (
	"context"
	"errors"

	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/errs"
	"github.com/permission-dev/internal/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const BizIDName = "biz_id"

type InterceptorBuilder struct {
	token *jwt.Token
}

func New(token *jwt.Token) *InterceptorBuilder {
	return &InterceptorBuilder{
		token: token,
	}
}

func (ib *InterceptorBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		//提取metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		//获取Authorization头
		authHeaders := md.Get("Authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization is required")
		}
		//处理token
		tokenStr := authHeaders[0]
		claim, err := ib.token.Decode(tokenStr)
		if err != nil {
			//细化错误返回
			if errors.Is(err, jwt.ErrTokenExpired) {
				return nil, status.Error(codes.Unauthenticated, "token expired")
			}
			if errors.Is(err, jwt.ErrTokenSignatureInvaild) {
				return nil, status.Error(codes.Unauthenticated, "invaild signatrue")
			}
			return nil, status.Error(codes.Unauthenticated, "invalid token"+err.Error())
		}
		val, ok := claim[BizIDName]

		if ok {
			bizId := val.(float64)
			ctx = context.WithValue(ctx, BizIDName, int64(bizId))
			elog.Info("用户请求信息", elog.FieldExtMessage(bizId))
		}
		return handler(ctx, req)
	}
}
func GetBizIDFromContext(ctx context.Context) (int64, error) {
	val := ctx.Value(BizIDName)
	if val == nil {
		return 0, errs.ErrBizIDNotFound
	}
	v, ok := val.(int64)
	if !ok {
		return 0, errs.ErrBizIDNotFound
	}
	return v, nil
}
