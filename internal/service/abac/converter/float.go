package converter

import "strconv"

type FloatConverter struct {
}

func (f *FloatConverter) Decode(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func (f *FloatConverter) Encode(t float64) (string, error) {
	return strconv.FormatFloat(t, 'f', -1, 64), nil
}

func NewFloatConverter() *FloatConverter {
	return &FloatConverter{}
}
