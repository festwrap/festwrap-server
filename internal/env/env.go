package env

import (
	"fmt"
	"os"
	"strconv"
)

type EnvValue interface {
	~int | ~string
}

func GetEnvWithDefault[T EnvValue](key string, defaultValue T) (T, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}

	var result T
	var err error
	switch valueType := any(defaultValue).(type) {
	case int:
		parsed, parseErr := strconv.Atoi(value)
		if parseErr != nil {
			err = fmt.Errorf("could not convert %s into integer", value)
		} else {
			result = any(parsed).(T)
		}
	case string:
		result = any(value).(T)
	default:
		err = fmt.Errorf("unsupported type %v", valueType)
	}

	return result, err
}
