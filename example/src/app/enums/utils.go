package enums

import "github.com/pkg/errors"

func valueToString(value any) (string, error) {
	s, ok := value.(string)
	if !ok {
		b, ok := value.([]byte)
		if !ok {
			return "", errors.Errorf("failed to parse value to string: %+v", value)
		}
		s = string(b)
	}
	return s, nil
}
