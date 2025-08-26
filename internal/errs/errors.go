package errs

import "github.com/pkg/errors"

var (
	ErrPermissionDuplicate     = errors.New("权限记录biz、key、action唯一索引冲突")
	ErrRoleDuplicate           = errors.New("角色记录biz、type、name唯一索引冲突")
	ErrRolePermissionDuplicate = errors.New("角色权限关联记录唯一索引冲突")
	ErrResourceDuplicate       = errors.New("资源记录biz、type、key唯一索引冲突")
	ErrSupportedSignAlgorithm  = errors.New("不支持的签名算法")
	ErrDecodeJWTTokenFailed    = errors.New("JWT令牌解析失败")
	ErrInvalidJWTToken         = errors.New("无效的令牌")
	ErrBizIDNotFound           = errors.New("BizID不存在")
	ErrUnkonwOperator          = errors.New("未知操作")
	ErrUnkonwDataType          = errors.New("未知类型")
)
