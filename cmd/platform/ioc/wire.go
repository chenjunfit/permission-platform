//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/permission-dev/internal/api/grpc/rbac"
	"github.com/permission-dev/internal/ioc"
	"github.com/permission-dev/internal/repository"
	"github.com/permission-dev/internal/repository/dao"
	rbacSvc "github.com/permission-dev/internal/service/rbac"
)

var (
	baseSet = wire.NewSet(
		ioc.InitDB,
		ioc.InitJwtToken,
	)
	rbacSet = wire.NewSet(
		dao.NewRoleDao,
		dao.NewResourceDao,
		dao.NewPermissionDAO,
		dao.NewUserDaoDAO,
		dao.NewRolePermissionDAO,
		dao.NewUserPermissionDAO,
		dao.NewRoleInclusionDAO,
		dao.NewBusinessConfigDAO,
		repository.NewRoleRepository,
		repository.NewResourceRepository,
		repository.NewPermissionRepository,
		repository.NewUserRoleRepository,
		repository.NewRolePermissionRepository,
		repository.NewUserPermissionRepository,
		repository.NewRoleIncludeRepository,
		repository.NewBusinessConfigRepository,

		rbacSvc.NewService,
		rbacSvc.NewPermissionService,
	)
)

func InitApp() *ioc.App {
	wire.Build(
		// 基础设施
		baseSet,

		// RBAC 服务
		rbacSet,

		// GRPC服务器
		rbac.NewServer,
		rbac.NewPermissionServer,
		ioc.InitGRPC,
		wire.Struct(new(ioc.App), "*"),
	)

	return new(ioc.App)
}
