package env

import (
	"os"
	"strconv"
	"time"
)

type EnvPair[T any] struct {
	Value     T
	IsDefault bool
}

func GetVariable[T any](variableName string, fnCast func(string) (T, error), defaultValue T) (*EnvPair[T], error) {
	val, exists := os.LookupEnv(variableName)
	if !exists {
		return &EnvPair[T]{defaultValue, true}, nil
	}
	castVal, err := fnCast(val)
	if err != nil {
		return nil, err
	}
	return &EnvPair[T]{castVal, false}, nil
}

func CastString(val string) (string, error) {
	return val, nil
}

func CastBool(val string) (bool, error) {
	bVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return bVal, nil
}

func CastDuration(val string) (time.Duration, error) {
	dVal, err := time.ParseDuration(val)
	if err != nil {
		return 0, err
	}
	return dVal, nil
}
