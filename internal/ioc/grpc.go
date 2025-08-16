package ioc

import (
	"github.com/gotomicro/ego/server/egrpc"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
	"github.com/permission-dev/internal/api/grpc/rbac"
	"github.com/permission-dev/internal/pkg/jwt"
)

func InitGRPC(
	crudServer *rbac.Server,
	permissionServer *rbac.PermissionServer,
	token *jwt.Token,
) []*egrpc.Component {
	authInterceptor := auth.New(token).Build()
	rbacServer := egrpc.Load("server.grpc.rbac").Build(
		egrpc.WithUnaryInterceptor(authInterceptor),
	)
	permissionv1.RegisterRBACServiceServer(rbacServer.Server, crudServer)
	permissionv1.RegisterPermissionServiceServer(rbacServer.Server, permissionServer)
	return []*egrpc.Component{rbacServer}
}
