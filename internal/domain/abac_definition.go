package domain

type AttributeDefinition struct {
	ID             int64
	Name           string
	Description    string
	DataType       DataType   //属性类型
	EntityType     EntityType //主体类型
	ValidationRule string
	Ctime          int64
	Utime          int64
}
type DataType string

func (d DataType) String() string {
	return string(d)
}

const (
	DataTypeString   DataType = "string"
	DataTypeNumber   DataType = "number"
	DataTypeBoolean  DataType = "boolean"
	DataTypeFloat    DataType = "float"
	DataTypeDatetime DataType = "datetime"
	DataTypeArray    DataType = "array"
)

type EntityType string

func (e EntityType) String() string {
	return string(e)
}

const (
	ResourceTypeEntity    EntityType = "resource"
	SubjectTypeEntity     EntityType = "subject"
	EnvironmentTypeEntity EntityType = "environment"
)

type Attributes struct {
	Subject     SubAttrs // 属性名 name => 属性值 value
	Resource    SubAttrs
	Environment SubAttrs
}
type SubAttrs map[string]string

func (s SubAttrs) SetKv(k, v string) SubAttrs {
	if s == nil {
		s = map[string]string{
			k: v,
		}
	} else {
		s[k] = v
	}
	return s
}
