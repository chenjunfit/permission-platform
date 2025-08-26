package domain

type PolicyStatusType string

const (
	PolicyStatusActive   PolicyStatusType = "active"
	PolicyStatusInActive PolicyStatusType = "inactive"

	LogicType ExecuteType = "logic" // 逻辑运算符执行方法
)

type ExecuteType string

type Policy struct {
	ID          int64
	BizID       int64
	Name        string
	Description string
	ExecuteType ExecuteType
	Status      PolicyStatusType
	Permissions []UserPermission
	Rules       []PolicyRule
	Ctime       int64
	Utime       int64
}

func (p Policy) ContainsAnyPermissions(permissionIDs []int64) bool {
	for idx := range permissionIDs {
		permissionID := permissionIDs[idx]
		for jdx := range p.Permissions {
			permission := p.Permissions[jdx]
			if permission.ID == permissionID {
				return true
			}
		}
	}
	return false
}

type PolicyRule struct {
	ID        int64
	AttrDef   AttributeDefinition
	Value     string
	LeftRule  *PolicyRule
	RightRule *PolicyRule
	Operator  RuleOperator
	Ctime     int64
	Utime     int64
}

func (p PolicyRule) SafeLeft() PolicyRule {
	if p.LeftRule == nil {
		return PolicyRule{}
	}
	return *p.LeftRule
}

func (p PolicyRule) SafeRight() PolicyRule {
	if p.RightRule == nil {
		return PolicyRule{}
	}
	return *p.RightRule
}
