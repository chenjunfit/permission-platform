package rbac

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	permissionv1 "github.com/permission-dev/api/proto/gen/permission/v1"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/service/rbac"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func NewServer(rbac rbac.Service) *Server {
	return &Server{
		rbacService: rbac,
	}
}

type Server struct {
	permissionv1.UnimplementedRBACServiceServer
	rbacService rbac.Service
	baseServer
}

func (s *Server) CreateBusinessConfig(ctx context.Context, in *permissionv1.CreateBusinessConfigRequest) (*permissionv1.CreateBusinessConfigResponse, error) {
	if in.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "业务配置不能为空")
	}

	// 将proto中的业务配置转换为领域模型
	in.Config.Id = 0
	domainConfig := s.toBusniessConfigDomain(in.Config)

	// 调用服务创建业务配置
	created, err := s.rbacService.CreateBusinessConfig(ctx, domainConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, "创建业务配置失败: "+err.Error())
	}

	// 将领域模型转换回proto
	return &permissionv1.CreateBusinessConfigResponse{
		Config: s.toBusniessConfigProto(created),
	}, nil
}

func (s *Server) GetBusinessConfig(ctx context.Context, in *permissionv1.GetBusinessConfigRequest) (*permissionv1.GetBusinessConfigResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务ID必须大于0")
	}
	// 调用服务获取业务配置
	config, err := s.rbacService.GetBusinessConfigByID(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取业务配置失败: "+err.Error())
	}
	// 将领域模型转换为proto响应
	return &permissionv1.GetBusinessConfigResponse{
		Config: s.toBusniessConfigProto(config),
	}, nil
}

func (s *Server) UpdateBusinessConfig(ctx context.Context, in *permissionv1.UpdateBusinessConfigRequest) (*permissionv1.UpdateBusinessConfigResponse, error) {
	if in.Config == nil || in.Config.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务配置不能为空且ID必须大于0")
	}

	// 将proto中的业务配置转换为领域模型
	domainConfig := s.toBusniessConfigDomain(in.Config)

	// 调用服务更新业务配置
	_, err := s.rbacService.UpdateBusinessConfig(ctx, domainConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, "更新业务配置失败: "+err.Error())
	}

	return &permissionv1.UpdateBusinessConfigResponse{
		Success: true,
	}, nil
}

func (s *Server) DeleteBusinessConfig(ctx context.Context, in *permissionv1.DeleteBusinessConfigRequest) (*permissionv1.DeleteBusinessConfigResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务ID必须大于0")
	}

	// 调用服务删除业务配置
	err := s.rbacService.DeleteBusinessConfigByID(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "删除业务配置失败: "+err.Error())
	}

	return &permissionv1.DeleteBusinessConfigResponse{
		Success: true,
	}, nil
}

func (s *Server) ListBusinessConfigs(ctx context.Context, in *permissionv1.ListBusinessConfigsRequest) (*permissionv1.ListBusinessConfigsResponse, error) {
	// 参数校验 - 设置默认分页
	offset := int(in.Offset)
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 10 // 默认每页10条
	}
	// 调用服务获取业务配置列表
	configs, err := s.rbacService.ListBusinessConfigs(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取业务配置列表失败: "+err.Error())
	}
	return &permissionv1.ListBusinessConfigsResponse{
		Configs: slice.Map(configs, func(_ int, src domain.BusinessConfig) *permissionv1.BusinessConfig {
			return s.toBusniessConfigProto(src)
		}),
	}, nil
}

func (s *Server) CreateRoleInclusion(ctx context.Context, in *permissionv1.CreateRoleInclusionRequest) (*permissionv1.CreateRoleInclusionResponse, error) {
	if in.RoleInclusion == nil {
		return nil, status.Error(codes.InvalidArgument, "角色包含关系不能为空")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 构建domain层的RoleInclusion对象
	in.RoleInclusion.Id = 0
	in.RoleInclusion.BizId = bizID
	domainRoleInclusion := s.toRoleInclusionDomain(in.RoleInclusion)

	// 调用服务创建角色包含关系
	created, err := s.rbacService.CreateRoleInclusion(ctx, domainRoleInclusion)
	if err != nil {
		return nil, status.Error(codes.Internal, "创建角色包含关系失败: "+err.Error())
	}

	// 转换回proto
	return &permissionv1.CreateRoleInclusionResponse{
		RoleInclusion: s.toRoleInclusionProto(created),
	}, nil
}

func (s *Server) GetRoleInclusion(ctx context.Context, in *permissionv1.GetRoleInclusionRequest) (*permissionv1.GetRoleInclusionResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "角色包含关系ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务获取角色包含关系
	roleInclusion, err := s.rbacService.GetRoleInclusion(ctx, bizID, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取角色包含关系失败: "+err.Error())
	}

	// 转换为proto响应
	return &permissionv1.GetRoleInclusionResponse{
		RoleInclusion: s.toRoleInclusionProto(roleInclusion),
	}, nil
}

func (s *Server) DeleteRoleInclusion(ctx context.Context, in *permissionv1.DeleteRoleInclusionRequest) (*permissionv1.DeleteRoleInclusionResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "角色包含关系ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务删除角色包含关系
	err = s.rbacService.DeleteRoleInclusion(ctx, bizID, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "删除角色包含关系失败: "+err.Error())
	}

	return &permissionv1.DeleteRoleInclusionResponse{
		Success: true,
	}, nil
}

func (s *Server) ListRoleInclusions(ctx context.Context, in *permissionv1.ListRoleInclusionsRequest) (*permissionv1.ListRoleInclusionsResponse, error) {
	// 参数校验
	if in.BizId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务ID必须大于0")
	}

	offset := int(in.Offset)
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 10 // 默认每页10条
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务获取角色包含关系列表
	roleInclusions, err := s.rbacService.ListRoleInclusions(ctx, bizID, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取角色包含关系列表失败: "+err.Error())
	}

	// 将领域模型列表转换为proto响应
	protoRoleInclusions := slice.Map(roleInclusions, func(_ int, src domain.RoleInclusion) *permissionv1.RoleInclusion {
		return s.toRoleInclusionProto(src)
	})

	return &permissionv1.ListRoleInclusionsResponse{
		RoleInclusions: protoRoleInclusions,
	}, nil
}

func (s *Server) GrantUserPermission(ctx context.Context, in *permissionv1.GrantUserPermissionRequest) (*permissionv1.GrantUserPermissionResponse, error) {
	if in.UserPermission == nil {
		return nil, status.Error(codes.InvalidArgument, "用户权限关系不能为空")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 设置业务ID
	in.UserPermission.Id = 0
	in.UserPermission.BizId = bizID

	// 构建domain层的UserPermission对象
	domainUserPermission := s.toUserPermissionDomain(in.UserPermission)

	// 调用服务授予用户权限
	created, err := s.rbacService.GrantUserPermission(ctx, domainUserPermission)
	if err != nil {
		return nil, status.Error(codes.Internal, "授予用户权限失败: "+err.Error())
	}

	// 转换回proto
	return &permissionv1.GrantUserPermissionResponse{
		UserPermission: s.toUserPermissionProto(created),
	}, nil
}

func (s *Server) RevokeUserPermission(ctx context.Context, in *permissionv1.RevokeUserPermissionRequest) (*permissionv1.RevokeUserPermissionResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户权限关系ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务撤销用户权限
	err = s.rbacService.RevokeUserPermission(ctx, bizID, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "撤销用户权限失败: "+err.Error())
	}

	return &permissionv1.RevokeUserPermissionResponse{
		Success: true,
	}, nil
}

func (s *Server) ListUserPermissions(ctx context.Context, in *permissionv1.ListUserPermissionsRequest) (*permissionv1.ListUserPermissionsResponse, error) {
	if in.BizId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务ID必须大于0")
	}

	offset := int(in.Offset)
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 10 // 默认每页10条
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务获取用户权限列表
	userPermissions, err := s.rbacService.ListUserPermissions(ctx, bizID, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取用户权限列表失败: "+err.Error())
	}

	// 将领域模型列表转换为proto响应
	return &permissionv1.ListUserPermissionsResponse{
		UserPermissions: slice.Map(userPermissions, func(_ int, src domain.UserPermission) *permissionv1.UserPermission {
			return s.toUserPermissionProto(src)
		}),
	}, nil
}

func (s *Server) GrantRolePermission(ctx context.Context, in *permissionv1.GrantRolePermissionRequest) (*permissionv1.GrantRolePermissionResponse, error) {
	if in.RolePermission == nil {
		return nil, status.Error(codes.InvalidArgument, "角色权限关系不能为空")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 构建domain层的RolePermission对象
	in.RolePermission.Id = 0
	in.RolePermission.BizId = bizID
	domainRolePermission := s.toRolePermissionDomain(in.RolePermission)

	// 调用服务授予角色权限
	created, err := s.rbacService.GrantRolePermission(ctx, domainRolePermission)
	if err != nil {
		return nil, status.Error(codes.Internal, "授予角色权限失败: "+err.Error())
	}

	// 转换回proto
	return &permissionv1.GrantRolePermissionResponse{
		RolePermission: s.toRolePermissionProto(created),
	}, nil
}

func (s *Server) RevokeRolePermission(ctx context.Context, in *permissionv1.RevokeRolePermissionRequest) (*permissionv1.RevokeRolePermissionResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "角色权限关系ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务撤销角色权限
	err = s.rbacService.RevokeRolePermission(ctx, bizID, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "撤销角色权限失败: "+err.Error())
	}

	return &permissionv1.RevokeRolePermissionResponse{
		Success: true,
	}, nil
}

func (s *Server) ListRolePermissions(ctx context.Context, in *permissionv1.ListRolePermissionsRequest) (*permissionv1.ListRolePermissionsResponse, error) {
	if in.BizId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务获取角色权限列表
	rolePermissions, err := s.rbacService.ListRolePermissions(ctx, bizID)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取角色权限列表失败: "+err.Error())
	}

	return &permissionv1.ListRolePermissionsResponse{
		RolePermissions: slice.Map(rolePermissions, func(_ int, src domain.RolePermission) *permissionv1.RolePermission {
			return s.toRolePermissionProto(src)
		}),
	}, nil
}

func (s *Server) GrantUserRole(ctx context.Context, in *permissionv1.GrantUserRoleRequest) (*permissionv1.GrantUserRoleResponse, error) {
	if in.UserRole == nil {
		return nil, status.Error(codes.InvalidArgument, "用户角色关系不能为空")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	in.UserRole.Id = 0
	in.UserRole.BizId = biz_id
	created, err := s.rbacService.GrantUserRole(ctx, s.toUserRoleDomain(in.UserRole))
	if err != nil {
		return nil, status.Error(codes.Internal, "授予用户角色失败"+err.Error())
	}
	return &permissionv1.GrantUserRoleResponse{UserRole: s.toUserRoleProto(created)}, nil
}

func (s *Server) RevokeUserRole(ctx context.Context, in *permissionv1.RevokeUserRoleRequest) (*permissionv1.RevokeUserRoleResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "用户角色关系ID必须大于0")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())

	}
	err = s.rbacService.RevokeUserRole(ctx, biz_id, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error()+"撤销用户角色失败")
	}
	return &permissionv1.RevokeUserRoleResponse{Success: true}, nil
}

func (s *Server) ListUserRoles(ctx context.Context, in *permissionv1.ListUserRolesRequest) (*permissionv1.ListUserRolesResponse, error) {
	if in.BizId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "BizID必须大于0")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())

	}
	userRoles, err := s.rbacService.ListUserRoles(ctx, biz_id)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询失败")
	}
	return &permissionv1.ListUserRolesResponse{UserRoles: slice.Map(userRoles, func(idx int, src domain.UserRole) *permissionv1.UserRole {
		return s.toUserRoleProto(src)
	})}, nil

}

func (s *Server) CreatePermission(ctx context.Context, in *permissionv1.CreatePermissionRequest) (*permissionv1.CreatePermissionResponse, error) {
	if in.Permission == nil {
		return nil, status.Error(codes.InvalidArgument, "权限不能为空")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	in.Permission.Id = 0
	in.Permission.BizId = biz_id
	created, err := s.rbacService.CreatePermission(ctx, s.toPermissionDomain(in.Permission))
	if err != nil {
		return nil, status.Error(codes.Internal, "创建权限失败")
	}
	return &permissionv1.CreatePermissionResponse{Permission: s.toPermissionProto(created)}, nil
}

func (s *Server) GetPermission(ctx context.Context, in *permissionv1.GetPermissionRequest) (*permissionv1.GetPermissionResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "权限不能为空")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	domainPermission, err := s.rbacService.GetPermission(ctx, biz_id, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询失败")
	}
	return &permissionv1.GetPermissionResponse{Permission: s.toPermissionProto(domainPermission)}, nil
}

func (s *Server) UpdatePermission(ctx context.Context, in *permissionv1.UpdatePermissionRequest) (*permissionv1.UpdatePermissionResponse, error) {
	if in.Permission == nil || in.Permission.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "权限不能为空且ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 构建domain层的Permission对象
	in.Permission.BizId = bizID
	domainPermission := s.toPermissionDomain(in.Permission)

	// 调用服务更新权限
	_, err = s.rbacService.UpdatePermission(ctx, domainPermission)
	if err != nil {
		return nil, status.Error(codes.Internal, "更新权限失败: "+err.Error())
	}

	return &permissionv1.UpdatePermissionResponse{Success: true}, nil
}

func (s *Server) DeletePermission(ctx context.Context, in *permissionv1.DeletePermissionRequest) (*permissionv1.DeletePermissionResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "权限ID必须大于0")
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务删除权限
	err = s.rbacService.DeletePermission(ctx, bizID, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "删除权限失败: "+err.Error())
	}

	return &permissionv1.DeletePermissionResponse{
		Success: true,
	}, nil
}

func (s *Server) ListPermissions(ctx context.Context, in *permissionv1.ListPermissionsRequest) (*permissionv1.ListPermissionsResponse, error) {
	// 参数校验
	if in.BizId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "业务ID必须大于0")
	}

	offset := int(in.Offset)
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 10 // 默认每页10条
	}

	bizID, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// 调用服务获取权限列表
	permissions, err := s.rbacService.ListPermissions(ctx, bizID, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取权限列表失败: "+err.Error())
	}

	return &permissionv1.ListPermissionsResponse{
		Permissions: slice.Map(permissions, func(_ int, src domain.Permission) *permissionv1.Permission {
			return s.toPermissionProto(src)
		}),
	}, nil
}

func (s *Server) CreateResource(ctx context.Context, in *permissionv1.CreateResourceRequest) (*permissionv1.CreateResourceResponse, error) {
	if in.Resource == nil {
		return nil, status.Error(codes.InvalidArgument, "资源不能为空")
	}
	in.Resource.Id = 0
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	in.Resource.BizId = biz_id
	created, err := s.rbacService.CreateResource(ctx, s.toResourceDomain(in.Resource))
	if err != nil {
		return nil, status.Error(codes.Internal, "创建资源失败"+err.Error())
	}
	return &permissionv1.CreateResourceResponse{Resource: s.toResourceProto(created)}, nil
}

func (s *Server) GetResource(ctx context.Context, in *permissionv1.GetResourceRequest) (*permissionv1.GetResourceResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id无效")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	resource, err := s.rbacService.GetResource(ctx, biz_id, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &permissionv1.GetResourceResponse{Resource: s.toResourceProto(resource)}, nil
}

func (s *Server) UpdateResource(ctx context.Context, in *permissionv1.UpdateResourceRequest) (*permissionv1.UpdateResourceResponse, error) {
	if in.Resource == nil || in.Resource.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "资源不能为空，且id必须大于0")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	in.Resource.BizId = biz_id
	_, err = s.rbacService.UpdateResource(ctx, s.toResourceDomain(in.Resource))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &permissionv1.UpdateResourceResponse{Success: true}, nil
}

func (s *Server) DeleteResource(ctx context.Context, in *permissionv1.DeleteResourceRequest) (*permissionv1.DeleteResourceResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id必须大于0")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	err = s.rbacService.DeleteResource(ctx, biz_id, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &permissionv1.DeleteResourceResponse{Success: true}, nil
}

func (s *Server) ListResources(ctx context.Context, in *permissionv1.ListResourcesRequest) (*permissionv1.ListResourcesResponse, error) {
	offset := int(in.Offset)
	limit := int(in.Limit)
	if limit <= 0 {
		limit = 10
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	resources, err := s.rbacService.ListResources(ctx, biz_id, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询资源列表失败")
	}
	return &permissionv1.ListResourcesResponse{
			Resources: slice.Map(resources, func(idx int, src domain.Resource) *permissionv1.Resource {
				return s.toResourceProto(src)
			}),
		},
		nil
}

func (s *Server) CreateRole(ctx context.Context, in *permissionv1.CreateRoleRequest) (*permissionv1.CreateRoleResponse, error) {
	if in.Role == nil {
		return nil, status.Error(codes.InvalidArgument, "role不能为空")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error()+"biz_id不能为空")
	}
	in.Role.Id = 0
	in.Role.BizId = biz_id
	domainRole := s.toRoleDomain(in.Role)
	created, err := s.rbacService.CreateRole(ctx, domainRole)
	if err != nil {
		return nil, status.Error(codes.Internal, "创建角色失败"+err.Error())
	}
	return &permissionv1.CreateRoleResponse{Role: s.toRoleProto(created)}, nil
}

func (s *Server) GetRole(ctx context.Context, in *permissionv1.GetRoleRequest) (*permissionv1.GetRoleResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "无效参数")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error()+"biz_id不能为空")
	}
	domainRole, err := s.rbacService.GetRole(ctx, biz_id, in.Id)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error()+"获取角色失败")
	}
	res := &permissionv1.GetRoleResponse{Role: s.toRoleProto(domainRole)}
	return res, nil
}

func (s *Server) UpdateRole(ctx context.Context, in *permissionv1.UpdateRoleRequest) (*permissionv1.UpdateRoleResponse, error) {
	if in.Role == nil || in.Role.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "无效参数")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error()+"biz_id不能为空")
	}
	in.Role.BizId = biz_id
	_, err = s.rbacService.UpdateRole(ctx, s.toRoleDomain(in.Role))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error()+"修改角色失败")
	}
	res := &permissionv1.UpdateRoleResponse{Success: true}
	return res, nil
}

func (s *Server) DeleteRole(ctx context.Context, in *permissionv1.DeleteRoleRequest) (*permissionv1.DeleteRoleResponse, error) {
	if in.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id无效")
	}
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "biz_id不能为空")
	}
	err = s.rbacService.DeleteRole(ctx, biz_id, in.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "删除角色失败")
	}
	return &permissionv1.DeleteRoleResponse{Success: true}, nil
}

func (s *Server) ListRoles(ctx context.Context, in *permissionv1.ListRolesRequest) (*permissionv1.ListRolesResponse, error) {
	biz_id, err := s.getBizIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "无效的biz_id")
	}
	limit := int(in.Limit)
	offset := int(in.Offset)
	if limit <= 0 {
		limit = 10 // 默认每页10条
	}
	roles, err := s.rbacService.ListRolesByRoleType(ctx, biz_id, in.Type, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error()+"获取角色列表失败")
	}
	return &permissionv1.ListRolesResponse{
		Roles: slice.Map(roles, func(idx int, src domain.Role) *permissionv1.Role {
			return s.toRoleProto(src)
		}),
	}, nil

}
func (s *Server) toRoleDomain(req *permissionv1.Role) domain.Role {
	var md string
	if req.Metadata != "" {
		md = req.Metadata
	}
	return domain.Role{
		ID:          req.Id,
		BizID:       req.BizId,
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		Metadata:    md,
	}
}
func (s *Server) toRoleProto(created domain.Role) *permissionv1.Role {
	return &permissionv1.Role{
		Id:          created.ID,
		BizId:       created.BizID,
		Type:        created.Type,
		Name:        created.Name,
		Description: created.Description,
		Metadata:    created.Metadata,
	}
}
func (s *Server) toResourceDomain(req *permissionv1.Resource) domain.Resource {
	var md string
	if req.Metadata != "" {
		md = req.Metadata
	}

	return domain.Resource{
		ID:          req.Id,
		BizID:       req.BizId,
		Type:        req.Type,
		Key:         req.Key,
		Name:        req.Name,
		Description: req.Description,
		Metadata:    md,
	}
}

func (s *Server) toResourceProto(created domain.Resource) *permissionv1.Resource {
	return &permissionv1.Resource{
		Id:          created.ID,
		BizId:       created.BizID,
		Type:        created.Type,
		Key:         created.Key,
		Name:        created.Name,
		Description: created.Description,
		Metadata:    created.Metadata,
	}
}
func (s *Server) toPermissionDomain(in *permissionv1.Permission) domain.Permission {
	var actions string
	if len(in.Actions) > 0 {
		actions = in.Actions[0]
	}
	var metadata string
	if in.Metadata != "" {
		metadata = in.Metadata
	}
	return domain.Permission{
		ID:          in.Id,
		BizID:       in.BizId,
		Name:        in.Name,
		Description: in.Description,
		Resource: domain.Resource{
			ID:   in.ResourceId,
			Type: in.ResourceType,
			Key:  in.ResourceKey,
		},
		Action:   actions,
		Metadata: metadata,
	}
}
func (s *Server) toPermissionProto(permission domain.Permission) *permissionv1.Permission {
	var actions []string
	if permission.Action != "" {
		actions = append(actions, permission.Action)
	}
	return &permissionv1.Permission{
		Id:           permission.ID,
		BizId:        permission.BizID,
		Name:         permission.Name,
		Description:  permission.Description,
		ResourceId:   permission.Resource.ID,
		ResourceType: permission.Resource.Type,
		ResourceKey:  permission.Resource.Key,
		Actions:      actions,
		Metadata:     permission.Metadata,
	}
}

func (s *Server) toUserRoleProto(userRole domain.UserRole) *permissionv1.UserRole {
	return &permissionv1.UserRole{
		Id:        userRole.ID,
		BizId:     userRole.BizID,
		UserId:    userRole.UserID,
		RoleId:    userRole.Role.ID,
		RoleName:  userRole.Role.Name,
		RoleType:  userRole.Role.Type,
		StartTime: userRole.StartTime,
		EndTime:   userRole.EndTime,
	}
}
func (s *Server) toUserRoleDomain(userRole *permissionv1.UserRole) domain.UserRole {
	start_time := userRole.StartTime
	end_time := userRole.EndTime
	if start_time == 0 {
		start_time = time.Now().Unix()
	}
	if end_time == 0 {
		end_time = time.Now().Unix()
	}
	return domain.UserRole{
		ID:     userRole.Id,
		BizID:  userRole.BizId,
		UserID: userRole.UserId,
		Role: domain.Role{
			ID:    userRole.RoleId,
			BizID: userRole.BizId,
			Type:  userRole.RoleType,
			Name:  userRole.RoleName,
		},
		StartTime: start_time,
		EndTime:   end_time,
	}
}

func (s *Server) toRolePermissionProto(rp domain.RolePermission) *permissionv1.RolePermission {
	return &permissionv1.RolePermission{
		Id:               rp.ID,
		BizId:            rp.BizID,
		RoleId:           rp.Role.ID,
		PermissionId:     rp.Permission.ID,
		RoleName:         rp.Role.Name,
		RoleType:         rp.Role.Type,
		ResourceType:     rp.Permission.Resource.Type,
		ResourceKey:      rp.Permission.Resource.Key,
		PermissionAction: rp.Permission.Action,
	}
}
func (s *Server) toRolePermissionDomain(rp *permissionv1.RolePermission) domain.RolePermission {
	return domain.RolePermission{
		ID:    rp.Id,
		BizID: rp.BizId,
		Role: domain.Role{
			ID:   rp.RoleId,
			Type: rp.RoleType,
			Name: rp.RoleName,
		},
		Permission: domain.Permission{
			ID: rp.PermissionId,
			Resource: domain.Resource{
				Type: rp.ResourceType,
				Key:  rp.ResourceKey,
			},
			Action: rp.PermissionAction,
		},
	}
}

func (s *Server) toUserPermissionDomain(up *permissionv1.UserPermission) domain.UserPermission {
	var effect domain.Effect
	if up.Effect == "allow" {
		effect = domain.EffectAllow
	} else {
		effect = domain.EffectDeny
	}
	start_time := up.StartTime
	end_time := up.EndTime
	if start_time == 0 {
		start_time = time.Now().Unix()
	}
	if end_time == 0 {
		end_time = time.Now().AddDate(100, 0, 0).Unix()
	}
	return domain.UserPermission{
		ID:     up.Id,
		BizID:  up.BizId,
		UserID: up.UserId,
		Permission: domain.Permission{
			ID:    up.PermissionId,
			BizID: up.BizId,
			Name:  up.PermissionName,
			Resource: domain.Resource{
				Type: up.ResourceType,
				Key:  up.ResourceKey,
			},
			Action: up.PermissionAction,
		},
		StartTime: start_time,
		EndTime:   end_time,
		Effect:    effect,
	}
}
func (s *Server) toUserPermissionProto(up domain.UserPermission) *permissionv1.UserPermission {
	return &permissionv1.UserPermission{
		Id:               up.ID,
		BizId:            up.BizID,
		UserId:           up.UserID,
		PermissionId:     up.Permission.ID,
		PermissionName:   up.Permission.Name,
		ResourceType:     up.Permission.Resource.Type,
		ResourceKey:      up.Permission.Resource.Key,
		PermissionAction: up.Permission.Action,
		Effect:           up.Effect.String(),
		StartTime:        up.StartTime,
		EndTime:          up.EndTime,
	}
}
func (s *Server) toRoleInclusionDomain(ri *permissionv1.RoleInclusion) domain.RoleInclusion {
	return domain.RoleInclusion{
		ID:    ri.Id,
		BizID: ri.BizId,
		IncludingRole: domain.Role{
			ID:   ri.IncludingRoleId,
			Type: ri.IncludingRoleType,
			Name: ri.IncludingRoleName,
		},
		IncludedRole: domain.Role{
			ID:   ri.IncludedRoleId,
			Type: ri.IncludedRoleType,
			Name: ri.IncludedRoleName,
		},
	}
}

func (s *Server) toRoleInclusionProto(created domain.RoleInclusion) *permissionv1.RoleInclusion {
	return &permissionv1.RoleInclusion{
		Id:                created.ID,
		BizId:             created.BizID,
		IncludingRoleId:   created.IncludingRole.ID,
		IncludingRoleType: created.IncludingRole.Type,
		IncludingRoleName: created.IncludingRole.Name,
		IncludedRoleId:    created.IncludedRole.ID,
		IncludedRoleType:  created.IncludedRole.Type,
		IncludedRoleName:  created.IncludedRole.Name,
	}
}

func (s *Server) toBusniessConfigProto(config domain.BusinessConfig) *permissionv1.BusinessConfig {
	return &permissionv1.BusinessConfig{
		Id:        config.ID,
		OwnerId:   config.OwnerID,
		OwnerType: config.OwnerType,
		Name:      config.Name,
		RateLimit: int32(config.RateLimit),
		Token:     config.Token,
		Ctime:     config.Ctime,
		Utime:     config.Utime,
	}
}
func (s *Server) toBusniessConfigDomain(config *permissionv1.BusinessConfig) domain.BusinessConfig {
	return domain.BusinessConfig{
		ID:        config.Id,
		OwnerID:   config.OwnerId,
		OwnerType: config.OwnerType,
		Name:      config.Name,
		RateLimit: int(config.RateLimit),
		Token:     config.Token,
		Ctime:     config.Ctime,
		Utime:     config.Utime,
	}
}
