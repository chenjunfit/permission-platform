//go:build wireinject

package rbac

import (
	"github.com/google/wire"
	"github.com/permission-dev/internal/repository"
	"github.com/permission-dev/internal/repository/dao"
	"github.com/permission-dev/internal/service/rbac"
	"github.com/permission-dev/internal/test/ioc"
)

type Service struct {
	RoleRepo           repository.RoleRepository
	ResourceRepo       repository.ResourceRepository
	PermissionRepo     repository.PermissionRepository
	UserRoleRepo       repository.UserRoleRepository
	RolePermissionRepo repository.RolePermissionRepository
	RoleIncludeRepo    repository.RoleIncludeRepository
	BusinessConfigRepo repository.BusinessConfigRepository
	UserPermissionRepo repository.UserPermissionRepository
	Svc                rbac.Service
}

func Init() *Service {
	wire.Build(
		ioc.BaseSet,
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
		repository.NewRoleIncludeRepository,
		repository.NewUserPermissionRepository,
		repository.NewBusinessConfigRepository,
		rbac.NewService,
		wire.Struct(new(Service), "*"),
	)
	return nil
}
