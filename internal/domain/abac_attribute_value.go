package domain

import "github.com/ecodeclub/ekit/slice"

type AttributeValue struct {
	ID      int64
	AttrDef AttributeDefinition
	Value   string
	Ctime   int64
	Utime   int64
}

type ABACObject struct {
	ID         int64
	BizId      int64
	AttrValues []AttributeValue
}

func (s *ABACObject) ValuesMap() map[int64]AttributeValue {
	return slice.ToMapV(s.AttrValues, func(element AttributeValue) (int64, AttributeValue) {
		return element.AttrDef.ID, element
	})
}

func (s *ABACObject) FillDefinitions(attrs AttrDefs) {
	for index := range s.AttrValues {
		attrValue := s.AttrValues[index]
		if attrDefinition, ok := attrs.GetByID(attrValue.ID); ok {
			attrValue.AttrDef = attrDefinition
		}
	}
}
func (s *ABACObject) MergeRealTimeAttr(attrs AttrDefs, values map[string]string) {
	for key, value := range values {
		if attrDefinition, ok := attrs.GetByName(key); ok {
			s.SetAttributeVal(value, attrDefinition)
		}
	}
}
func (s *ABACObject) SetAttributeVal(val string, definition AttributeDefinition) {
	for idx := range s.AttrValues {
		if s.AttrValues[idx].AttrDef.ID == definition.ID {
			s.AttrValues[idx].Value = val
			return
		}
	}
	s.AttrValues = append(s.AttrValues, AttributeValue{
		AttrDef: definition,
		Value:   val,
	})
}
