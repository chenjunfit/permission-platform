package domain

const (
	Equals         RuleOperator = "="
	NotEquals      RuleOperator = "!="
	Greater        RuleOperator = ">"
	Less           RuleOperator = "<"
	GreaterOrEqual RuleOperator = ">="
	LessOrEqual    RuleOperator = "<="
	AND            RuleOperator = "AND"
	OR             RuleOperator = "OR"
	IN             RuleOperator = "IN"
	NotIn          RuleOperator = "NOT IN"
	NOT            RuleOperator = "NOT"
	AllMatch       RuleOperator = "ALL MATCH"
	AnyMatch       RuleOperator = "ANY MATCH"
)

type RuleOperator string

func (r RuleOperator) String() string {
	return string(r)
}
