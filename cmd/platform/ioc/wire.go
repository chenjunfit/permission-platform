//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	"github.com/permission-dev/internal/api/grpc/rbac"
	"github.com/permission-dev/internal/ioc"
	"github.com/permission-dev/internal/repository"
	"github.com/permission-dev/internal/repository/dao"
	"github.com/permission-dev/internal/repository/dao/audit"
	"github.com/permission-dev/internal/service/abac"
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

		dao.NewAttributeDefinitionDAO,
		dao.NewResourceAttributeValueDAO,
		dao.NewEnvironmentAttributeValueDAO,
		dao.NewSubjectAttributeValueDAO,

		repository.NewRoleRepository,
		repository.NewResourceRepository,
		repository.NewPermissionRepository,
		repository.NewUserRoleRepository,
		repository.NewRolePermissionRepository,
		repository.NewUserPermissionRepository,
		repository.NewRoleIncludeRepository,
		repository.NewBusinessConfigRepository,

		repository.NewAttributeDefinitionRepository,
		repository.NewAttributeValueRepository,
		repository.NewAttributePolicyRepository,

		rbacSvc.NewService,
		rbacSvc.NewPermissionService,

		abac.NewAttributeDefinitionSvc,
		abac.NewAttributeValueSvc,
		abac.NewPolicySvc,

		audit.NewOperationLogDao,
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
