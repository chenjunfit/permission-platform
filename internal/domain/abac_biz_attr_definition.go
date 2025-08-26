package domain

type AttrDefs []AttributeDefinition

func (a AttrDefs) Map() map[int64]AttributeDefinition {
	res := make(map[int64]AttributeDefinition, len(a))
	for idx := range a {
		val := a[idx]
		res[val.ID] = val
	}
	return res
}
func (a AttrDefs) GetByID(id int64) (AttributeDefinition, bool) {
	for idx := range a {
		if a[idx].ID == id {
			return a[idx], true
		}
	}
	return AttributeDefinition{}, false
}

func (a AttrDefs) GetByName(name string) (AttributeDefinition, bool) {
	for idx := range a {
		if a[idx].Name == name {
			return a[idx], true
		}
	}
	return AttributeDefinition{}, false
}

type BizAttrDefinition struct {
	BizID               int64
	SubjectAttrDefs     AttrDefs
	ResourceAttrDefs    AttrDefs
	EnvironmentAttrDefs AttrDefs
	AllDefs             map[int64]AttributeDefinition
}

func (biz *BizAttrDefinition) GetByDefId(id int64) (AttributeDefinition, bool) {
	val, ok := biz.AllDefs[id]
	return val, ok
}
