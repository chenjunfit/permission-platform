package dao

import (
	"github.com/ego-component/egorm"
	"github.com/go-sql-driver/mysql"
	"github.com/permission-dev/internal/repository/dao/audit"
	"github.com/pkg/errors"
)

func InitTable(db *egorm.Component) error {
	return db.AutoMigrate(
		&Role{},
		&Resource{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
		&RoleInclusion{},
		&UserPermission{},
		&BusinessConfig{},

		&AttributeDefinition{},
		&EnvironmentAttributeValue{},
		&ResourceAttributeValue{},
		&SubjectAttributeValue{},
		&Policy{},
		&PolicyRule{},
		&PermissionPolicy{},

		&audit.OperationLog{},
		&audit.UserRoleLog{},
	)
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	me := new(mysql.MySQLError)
	if ok := errors.As(err, &me); ok {
		const uniqueIndexErrNo uint16 = 1062
		return me.Number == uniqueIndexErrNo
	}
	return false
}
