package whitelist

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"sync"
)

// 必须放在jwt后面,获取biz id
type InteceptorBuilder struct {
	//从配置文件读取的话，在init grpc中加载配置文件
	WhiteList []int64
	mutex     *sync.RWMutex
}

func NewInteceptorBuilder(whiteList []int64) *InteceptorBuilder {
	return &InteceptorBuilder{
		WhiteList: whiteList,
		mutex:     &sync.RWMutex{},
	}
}

func (i *InteceptorBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		i.mutex.RLock()
		defer i.mutex.RUnlock()
		if strings.Contains(info.FullMethod, "BusinessConfig") {
			bizId, err := auth.GetBizIDFromContext(ctx)
			if err != nil {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
			if !slice.Contains(i.WhiteList, bizId) {
				return nil, status.Errorf(codes.Unauthenticated, "不在白名单")
			}
		}
		return handler(ctx, req)
	}
}
func (i *InteceptorBuilder) UpdateWhiteList(list []int64) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.WhiteList = list
}
