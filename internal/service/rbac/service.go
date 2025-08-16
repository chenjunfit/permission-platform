package rbac

import (
	"context"
	"github.com/permission-dev/internal/api/grpc/interceptor/auth"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/pkg/jwt"
	"github.com/permission-dev/internal/repository"
	"time"
)

type Service interface {
	// 角色相关方法
	CreateRole(ctx context.Context, role domain.Role) (domain.Role, error)
	GetRole(ctx context.Context, bizID, id int64) (domain.Role, error)
	UpdateRole(ctx context.Context, role domain.Role) (domain.Role, error)
	DeleteRole(ctx context.Context, bizID, id int64) error
	ListRolesByRoleType(ctx context.Context, bizID int64, roleType string, offset, limit int) ([]domain.Role, error)
	ListRoles(ctx context.Context, bizID int64, offset, limit int) ([]domain.Role, error)
	//资源相关方法
	CreateResource(ctx context.Context, resource domain.Resource) (domain.Resource, error)
	GetResource(ctx context.Context, bizID, id int64) (domain.Resource, error)
	UpdateResource(ctx context.Context, resource domain.Resource) (domain.Resource, error)
	DeleteResource(ctx context.Context, bizID, id int64) error
	ListResources(ctx context.Context, bizID int64, offset, limit int) ([]domain.Resource, error)
	//权限相关方法
	CreatePermission(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	GetPermission(ctx context.Context, bizID, id int64) (domain.Permission, error)
	UpdatePermission(ctx context.Context, permission domain.Permission) (domain.Permission, error)
	DeletePermission(ctx context.Context, bizID, id int64) error
	ListPermissions(ctx context.Context, bizID int64, offset, limit int) ([]domain.Permission, error)
	//用户角色相关方法
	GrantUserRole(ctx context.Context, userRole domain.UserRole) (domain.UserRole, error)
	RevokeUserRole(ctx context.Context, bizID, id int64) error
	ListUserRolesByUserID(ctx context.Context, bizID, userID int64) ([]domain.UserRole, error)
	ListUserRoles(ctx context.Context, bizID int64) ([]domain.UserRole, error)
	//角色权限相关方法
	GrantRolePermission(ctx context.Context, rolePermission domain.RolePermission) (domain.RolePermission, error)
	RevokeRolePermission(ctx context.Context, bizID, id int64) error
	ListRolePermissionsByRoleID(ctx context.Context, bizID, roleID int64) ([]domain.RolePermission, error)
	ListRolePermissions(ctx context.Context, bizID int64) ([]domain.RolePermission, error)
	//角色包含相关方法
	CreateRoleInclusion(ctx context.Context, roleInclusion domain.RoleInclusion) (domain.RoleInclusion, error)
	GetRoleInclusion(ctx context.Context, bizID, id int64) (domain.RoleInclusion, error)
	DeleteRoleInclusion(ctx context.Context, bizID, id int64) error
	ListRoleInclusionsByRoleID(ctx context.Context, bizID, roleID int64, isIncluding bool) ([]domain.RoleInclusion, error)
	ListRoleInclusions(ctx context.Context, bizID int64, offset, limit int) ([]domain.RoleInclusion, error)
	//用户权限相关方法
	GrantUserPermission(ctx context.Context, userPermission domain.UserPermission) (domain.UserPermission, error)
	RevokeUserPermission(ctx context.Context, bizID, id int64) error
	ListUserPermissionsByUserID(ctx context.Context, bizID, userID int64) ([]domain.UserPermission, error)
	ListUserPermissions(ctx context.Context, bizID int64, offset, limit int) ([]domain.UserPermission, error)
	//业务接入相关方法
	CreateBusinessConfig(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error)
	GetBusinessConfigByID(ctx context.Context, id int64) (domain.BusinessConfig, error)
	UpdateBusinessConfig(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error)
	DeleteBusinessConfigByID(ctx context.Context, id int64) error
	ListBusinessConfigs(ctx context.Context, offset, limit int) ([]domain.BusinessConfig, error)
}

func NewService(
	roleRepo repository.RoleRepository,
	resourceRepo repository.ResourceRepository,
	permissionRepo repository.PermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	rolePermissionRepository repository.RolePermissionRepository,
	includeRepository repository.RoleIncludeRepository,
	userPermissionRepository repository.UserPermissionRepository,
	businessConfigRepository repository.BusinessConfigRepository,
	jwtToken *jwt.Token,
) Service {
	return &rbacService{
		roleRepo:                 roleRepo,
		resourceRepo:             resourceRepo,
		permissionRepo:           permissionRepo,
		userRoleRepo:             userRoleRepo,
		rolePermissionRepo:       rolePermissionRepository,
		roleIncludeRepo:          includeRepository,
		userPermissionRepo:       userPermissionRepository,
		businessConfigRepository: businessConfigRepository,
		jwtToken:                 jwtToken,
	}
}

type rbacService struct {
	roleRepo                 repository.RoleRepository
	resourceRepo             repository.ResourceRepository
	permissionRepo           repository.PermissionRepository
	userRoleRepo             repository.UserRoleRepository
	rolePermissionRepo       repository.RolePermissionRepository
	roleIncludeRepo          repository.RoleIncludeRepository
	userPermissionRepo       repository.UserPermissionRepository
	businessConfigRepository repository.BusinessConfigRepository
	jwtToken                 *jwt.Token
	
}

func (r *rbacService) CreateBusinessConfig(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error) {
	created, err := r.businessConfigRepository.Create(ctx, config)
	if err != nil {
		return domain.BusinessConfig{}, err
	}
	token, err := r.jwtToken.Encode(jwt.MapClaims{
		auth.BizIDName: created.ID,
		"exp":          time.Now().AddDate(100, 0, 0).Unix(),
	})
	if err != nil {
		return domain.BusinessConfig{}, err
	}
	created.Token = token
	err = r.businessConfigRepository.UpdateToken(ctx, created.ID, token)
	if err != nil {
		return domain.BusinessConfig{}, err
	}
	return created, nil
}

func (r *rbacService) GetBusinessConfigByID(ctx context.Context, id int64) (domain.BusinessConfig, error) {
	return r.businessConfigRepository.FindByID(ctx, id)
}

func (r *rbacService) UpdateBusinessConfig(ctx context.Context, config domain.BusinessConfig) (domain.BusinessConfig, error) {
	return r.UpdateBusinessConfig(ctx, config)
}

func (r *rbacService) DeleteBusinessConfigByID(ctx context.Context, id int64) error {
	return r.businessConfigRepository.Delete(ctx, id)
}

func (r *rbacService) ListBusinessConfigs(ctx context.Context, offset, limit int) ([]domain.BusinessConfig, error) {
	return r.businessConfigRepository.Find(ctx, offset, limit)
}

func (r *rbacService) GrantUserPermission(ctx context.Context, userPermission domain.UserPermission) (domain.UserPermission, error) {
	return r.userPermissionRepo.Create(ctx, userPermission)
}

func (r *rbacService) RevokeUserPermission(ctx context.Context, bizID, id int64) error {
	return r.userPermissionRepo.DeleteByBizIdAndID(ctx, bizID, id)
}

func (r *rbacService) ListUserPermissionsByUserID(ctx context.Context, bizID, userID int64) ([]domain.UserPermission, error) {
	return r.userPermissionRepo.FindByBizIdAndUserID(ctx, bizID, userID)
}

func (r *rbacService) ListUserPermissions(ctx context.Context, bizID int64, offset, limit int) ([]domain.UserPermission, error) {
	return r.userPermissionRepo.FindByBizID(ctx, bizID, offset, limit)
}

func (r *rbacService) CreateRoleInclusion(ctx context.Context, roleInclusion domain.RoleInclusion) (domain.RoleInclusion, error) {
	return r.roleIncludeRepo.Create(ctx, roleInclusion)
}

func (r *rbacService) GetRoleInclusion(ctx context.Context, bizID, id int64) (domain.RoleInclusion, error) {
	return r.roleIncludeRepo.FindByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) DeleteRoleInclusion(ctx context.Context, bizID, id int64) error {
	return r.roleIncludeRepo.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) ListRoleInclusionsByRoleID(ctx context.Context, bizID, roleID int64, isIncluding bool) ([]domain.RoleInclusion, error) {
	if isIncluding {
		return r.roleIncludeRepo.FindByBizIdAndIncludingIds(ctx, bizID, []int64{roleID})
	} else {
		return r.roleIncludeRepo.FindByBizIdAndIncludedIds(ctx, bizID, []int64{roleID})
	}
}

func (r *rbacService) ListRoleInclusions(ctx context.Context, bizID int64, offset, limit int) ([]domain.RoleInclusion, error) {
	return r.roleIncludeRepo.FindByBizID(ctx, bizID, offset, limit)
}

func (r *rbacService) GrantRolePermission(ctx context.Context, rolePermission domain.RolePermission) (domain.RolePermission, error) {
	return r.rolePermissionRepo.Create(ctx, rolePermission)
}

func (r *rbacService) RevokeRolePermission(ctx context.Context, bizID, id int64) error {
	return r.rolePermissionRepo.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) ListRolePermissionsByRoleID(ctx context.Context, bizID, roleID int64) ([]domain.RolePermission, error) {
	return r.rolePermissionRepo.FindByBizIDAndRoleIDs(ctx, bizID, []int64{roleID})
}

func (r *rbacService) ListRolePermissions(ctx context.Context, bizID int64) ([]domain.RolePermission, error) {
	return r.rolePermissionRepo.FindByBizID(ctx, bizID)
}

func (r *rbacService) GrantUserRole(ctx context.Context, userRole domain.UserRole) (domain.UserRole, error) {
	return r.userRoleRepo.Create(ctx, userRole)
}

func (r *rbacService) RevokeUserRole(ctx context.Context, bizID, id int64) error {
	return r.userRoleRepo.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) ListUserRolesByUserID(ctx context.Context, bizID, userID int64) ([]domain.UserRole, error) {
	return r.userRoleRepo.FindByBizIDAndUserID(ctx, bizID, userID)
}

func (r *rbacService) ListUserRoles(ctx context.Context, bizID int64) ([]domain.UserRole, error) {
	return r.userRoleRepo.FindByBizID(ctx, bizID)
}

func (r *rbacService) CreatePermission(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	return r.permissionRepo.Create(ctx, permission)
}

func (r *rbacService) GetPermission(ctx context.Context, bizID, id int64) (domain.Permission, error) {
	return r.permissionRepo.FindByBizIDANdID(ctx, bizID, id)
}

func (r *rbacService) UpdatePermission(ctx context.Context, permission domain.Permission) (domain.Permission, error) {
	return r.permissionRepo.UpdateByBizIDAndID(ctx, permission)
}

func (r *rbacService) DeletePermission(ctx context.Context, bizID, id int64) error {
	return r.permissionRepo.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) ListPermissions(ctx context.Context, bizID int64, offset, limit int) ([]domain.Permission, error) {
	return r.permissionRepo.FindByBizID(ctx, bizID, offset, limit)
}

func (r *rbacService) CreateResource(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	return r.resourceRepo.Create(ctx, resource)
}

func (r *rbacService) GetResource(ctx context.Context, bizID, id int64) (domain.Resource, error) {
	return r.resourceRepo.FindByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) UpdateResource(ctx context.Context, resource domain.Resource) (domain.Resource, error) {
	return r.resourceRepo.UpdateByBizIDAndID(ctx, resource)
}

func (r *rbacService) DeleteResource(ctx context.Context, bizID, id int64) error {
	return r.resourceRepo.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) ListResources(ctx context.Context, bizID int64, offset, limit int) ([]domain.Resource, error) {
	return r.resourceRepo.FindByBizID(ctx, bizID, offset, limit)
}

func (r *rbacService) CreateRole(ctx context.Context, role domain.Role) (domain.Role, error) {
	return r.roleRepo.Create(ctx, role)
}

func (r *rbacService) GetRole(ctx context.Context, bizID, id int64) (domain.Role, error) {
	return r.roleRepo.FindByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) UpdateRole(ctx context.Context, role domain.Role) (domain.Role, error) {
	return r.roleRepo.UpdateByBizIDAndID(ctx, role)
}

func (r *rbacService) DeleteRole(ctx context.Context, bizID, id int64) error {
	return r.roleRepo.DeleteByBizIDAndID(ctx, bizID, id)
}

func (r *rbacService) ListRolesByRoleType(ctx context.Context, bizID int64, roleType string, offset, limit int) ([]domain.Role, error) {
	return r.roleRepo.FindByBizIDAndType(ctx, bizID, roleType, offset, limit)
}

func (r *rbacService) ListRoles(ctx context.Context, bizID int64, offset, limit int) ([]domain.Role, error) {
	return r.roleRepo.FindByBizID(ctx, bizID, offset, limit)
}
