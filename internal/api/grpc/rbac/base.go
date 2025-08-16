package rbac

import (
	"context"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
)

type baseServer struct {
}

func (b *baseServer) getBizIDFromContext(ctx context.Context) (int64, error) {
	return auth.GetBizIDFromContext(ctx)
}
