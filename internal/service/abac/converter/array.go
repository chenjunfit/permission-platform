package converter

import "encoding/json"

type ArrayConverter struct {
}

func (a *ArrayConverter) Decode(str string) ([]string, error) {
	var res []string
	err := json.Unmarshal([]byte(str), &res)
	return res, err
}

func (a *ArrayConverter) Encode(t []string) (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func NewArrayConverter() *ArrayConverter {
	return &ArrayConverter{}
}
