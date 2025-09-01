package audit

import (
	"context"
	"encoding/json"
	"github.com/gotomicro/ego/core/elog"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
	"github.com/permission-dev/internal/repository/dao/audit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type InterceptorBuilder struct {
	dao    audit.OperationLogDAO
	logger *elog.Component
}

func NewInterceptorBuilder(dao audit.OperationLogDAO) *InterceptorBuilder {
	return &InterceptorBuilder{
		dao:    dao,
		logger: elog.DefaultLogger.With(elog.FieldName("audit.OperationLog")),
	}
}

func (b *InterceptorBuilder) Build() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var operationLog audit.OperationLog
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			operator := md.Get("operator")
			if len(operator) != 0 {
				operationLog.Operator = operator[0]
			}
			key := md.Get("key")
			if len(key) != 0 {
				operationLog.Key = key[0]
			}
		}
		operationLog.BizID, _ = auth.GetBizIDFromContext(ctx)
		operationLog.Method = info.FullMethod
		data, _ := json.Marshal(req)
		operationLog.Request = string(data)
		_, err = b.dao.Create(ctx, operationLog)
		if err != nil {
			b.logger.Error(
				"存储操作日志失败",
				elog.FieldErr(err),
				elog.FieldKey("operationLog"),
				elog.FieldValueAny(operationLog),
			)
		}
		return handler(ctx, req)
	}

}
