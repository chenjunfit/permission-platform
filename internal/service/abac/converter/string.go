package converter

type StringConverter struct {
}

func (s *StringConverter) Decode(str string) (string, error) {
	return str, nil
}

func (s *StringConverter) Encode(t string) (string, error) {
	return t, nil
}

func NewStringConverter() *StringConverter {
	return &StringConverter{}
}
