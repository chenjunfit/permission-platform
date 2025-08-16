package rbac

import (
	"context"
	"github.com/permission-dev/internal/domain"
	"github.com/permission-dev/internal/test/integration/ioc/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func cleanTestEnviroment(t *testing.T, ctx context.Context, svc interface{}) {
	t.Log("清理测试环境,保留预设数据")
	rbacSvc, ok := svc.(*rbac.Service)
	if !ok {
		t.Error("无法获取RBAC服务实例")
		return
	}
	const (
		safeBizID  = 2
		safeUserID = 1000
	)
	//清理角色包含关系
	roles, err := rbacSvc.Svc.ListRoles(ctx, safeBizID, 0, 1000)
	if err != nil {
		t.Errorf("无法获取biz_Id:%d下所有角色失败", safeUserID)
		return
	}
	for _, role := range roles {
		roleInclusions, err := rbacSvc.Svc.ListRoleInclusionsByRoleID(ctx, safeBizID, role.ID, true)
		if err != nil {
			t.Errorf("无法获取biz_Id:%d下角色包含关系失败", safeUserID)
			return
		}
		for _, roleInclusion := range roleInclusions {
			err := rbacSvc.Svc.DeleteRoleInclusion(ctx, safeBizID, roleInclusion.ID)
			if err != nil {
				t.Errorf("删除biz_Id:%d下角色包含关系:%d失败", safeUserID, roleInclusion.ID)
				return
			}
		}
	}

	//确保所有业务角色包含关系都被清理
	otherBizConfigs, err := rbacSvc.Svc.ListBusinessConfigs(ctx, 0, 1000)
	if err == nil {
		for _, config := range otherBizConfigs {
			//只清理非预设业务
			if config.ID > 1 {
				bizRoles, err := rbacSvc.Svc.ListRoles(ctx, config.ID, 0, 1000)
				if err == nil {
					// 对每个角色，清理其作为包含者的包含关系
					for _, role := range bizRoles {
						inclusions, err := rbacSvc.Svc.ListRoleInclusionsByRoleID(ctx, config.ID, role.ID, true)
						if err == nil {
							for _, ri := range inclusions {
								_ = rbacSvc.Svc.DeleteRoleInclusion(ctx, ri.BizID, ri.ID)
							}
						}
						inclusions, err = rbacSvc.Svc.ListRoleInclusionsByRoleID(ctx, config.ID, role.ID, false)
						if err == nil {
							for _, ri := range inclusions {
								_ = rbacSvc.Svc.DeleteRoleInclusion(ctx, ri.BizID, ri.ID)
							}
						}
					}
				}

			}

		}
	}

	// 1. 清理角色权限数据
	rolePermissions, err := rbacSvc.Svc.ListRolePermissions(ctx, safeBizID)
	if err == nil {
		for _, rp := range rolePermissions {
			err = rbacSvc.Svc.RevokeRolePermission(ctx, rp.BizID, rp.ID)
		}
	}

	// 2. 清理用户权限数据
	userPermissions, err := rbacSvc.Svc.ListUserPermissionsByUserID(ctx, safeBizID, safeUserID)
	if err == nil {
		for _, up := range userPermissions {
			_ = rbacSvc.Svc.RevokeUserPermission(ctx, up.BizID, up.ID)
		}
	}

	// 3. 清理用户角色数据
	userRoles, err := rbacSvc.Svc.ListUserRolesByUserID(ctx, safeBizID, safeUserID)
	if err == nil {
		for _, ur := range userRoles {
			_ = rbacSvc.Svc.RevokeUserRole(ctx, ur.BizID, ur.ID)
		}
	}

	// 4. 清理权限数据
	permissions, err := rbacSvc.Svc.ListPermissions(ctx, safeBizID, 0, 1000)
	if err == nil {
		for _, perm := range permissions {
			_ = rbacSvc.Svc.DeletePermission(ctx, perm.BizID, perm.ID)
		}
	}

	// 5. 清理角色数据
	roles, err = rbacSvc.Svc.ListRoles(ctx, safeBizID, 0, 1000)
	if err == nil {
		for _, role := range roles {
			err = rbacSvc.Svc.DeleteRole(ctx, role.BizID, role.ID)
		}
	}

	// 6. 清理资源数据
	resources, err := rbacSvc.Svc.ListResources(ctx, safeBizID, 0, 1000)
	if err == nil {
		for _, res := range resources {
			_ = rbacSvc.Svc.DeleteResource(ctx, res.BizID, res.ID)
		}
	}

	// 7. 清理每个测试业务的所有数据
	configs, err := rbacSvc.Svc.ListBusinessConfigs(ctx, 0, 1000)
	if err == nil {
		for _, config := range configs {
			// 只清理非预设数据 (bizID > 1)
			if config.ID > 1 {
				// 先删除该业务下的所有数据
				bizSpecificPermissions, err := rbacSvc.Svc.ListPermissions(ctx, config.ID, 0, 1000)
				if err == nil {
					for _, perm := range bizSpecificPermissions {
						err = rbacSvc.Svc.DeletePermission(ctx, perm.BizID, perm.ID)
					}
				}

				bizSpecificResources, err := rbacSvc.Svc.ListResources(ctx, config.ID, 0, 1000)
				if err == nil {
					for _, res := range bizSpecificResources {
						err = rbacSvc.Svc.DeleteResource(ctx, res.BizID, res.ID)
					}
				}

				bizSpecificRoles, err := rbacSvc.Svc.ListRoles(ctx, config.ID, 0, 1000)
				if err == nil {
					for _, role := range bizSpecificRoles {
						err = rbacSvc.Svc.DeleteRole(ctx, role.BizID, role.ID)
					}
				}

				// 最后删除业务配置
				err = rbacSvc.Svc.DeleteBusinessConfigByID(ctx, config.ID)

			}
		}
	}

}
func assertRole(t *testing.T, expected, actual domain.Role) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	// 不检查Metadata, Ctime和Utime，因为这些是系统自动生成的
}

// assertResource 断言资源对象的字段是否符合预期
func assertResource(t *testing.T, expected, actual domain.Resource) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	// 不检查Metadata, Ctime和Utime
}

// assertPermission 断言权限对象的字段是否符合预期
func assertPermission(t *testing.T, expected, actual domain.Permission) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description)
	if expected.Resource.ID > 0 {
		assert.Equal(t, expected.Resource.ID, actual.Resource.ID)
	}
	assert.Equal(t, expected.Resource.Type, actual.Resource.Type)
	assert.Equal(t, expected.Resource.Key, actual.Resource.Key)
	assert.Equal(t, expected.Action, actual.Action)
	// 不检查Metadata, Ctime和Utime
}

// assertUserRole 断言用户角色对象的字段是否符合预期
func assertUserRole(t *testing.T, expected, actual domain.UserRole) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	assert.Equal(t, expected.UserID, actual.UserID)
	if expected.Role.ID > 0 {
		assert.Equal(t, expected.Role.ID, actual.Role.ID)
	}
	assert.Equal(t, expected.StartTime, actual.StartTime)
	assert.Equal(t, expected.EndTime, actual.EndTime)
	// 不检查Ctime和Utime
}

// assertUserPermission 断言用户权限对象的字段是否符合预期
func assertUserPermission(t *testing.T, expected, actual domain.UserPermission) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	assert.Equal(t, expected.UserID, actual.UserID)
	if expected.Permission.ID > 0 {
		assert.Equal(t, expected.Permission.ID, actual.Permission.ID)
	}
	assert.Equal(t, expected.Effect, actual.Effect)
	assert.Equal(t, expected.StartTime, actual.StartTime)
	assert.Equal(t, expected.EndTime, actual.EndTime)
	// 不检查Ctime和Utime
}

// assertRolePermission 断言角色权限对象的字段是否符合预期
func assertRolePermission(t *testing.T, expected, actual domain.RolePermission) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	if expected.Role.ID > 0 {
		assert.Equal(t, expected.Role.ID, actual.Role.ID)
	}
	if expected.Permission.ID > 0 {
		assert.Equal(t, expected.Permission.ID, actual.Permission.ID)
	}
	// 不检查Ctime和Utime
}

// assertRoleInclusion 断言角色包含关系对象的字段是否符合预期
func assertRoleInclusion(t *testing.T, expected, actual domain.RoleInclusion) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.BizID, actual.BizID)
	if expected.IncludingRole.ID > 0 {
		assert.Equal(t, expected.IncludingRole.ID, actual.IncludingRole.ID)
	}
	if expected.IncludedRole.ID > 0 {
		assert.Equal(t, expected.IncludedRole.ID, actual.IncludedRole.ID)
	}
	// 不检查Ctime和Utime
}

// assertBusinessConfig 断言业务配置对象的字段是否符合预期
func assertBusinessConfig(t *testing.T, expected, actual domain.BusinessConfig) {
	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.OwnerID, actual.OwnerID)
	assert.Equal(t, expected.OwnerType, actual.OwnerType)
	assert.Equal(t, expected.RateLimit, actual.RateLimit)
	if expected.Token != "" {
		assert.NotEmpty(t, actual.Token)
	}
	// 不检查Ctime和Utime
}

// requireSetupCreate 要求创建操作成功
func requireSetupCreate(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
}

// cleanupTest 清理测试数据
func cleanupTest(t *testing.T, ctx context.Context, fn func() error) {
	if err := fn(); err != nil {
		t.Logf("清理测试数据失败: %v", err)
	}
}
