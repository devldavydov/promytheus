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

func GetVariable[T any](variableName string, fnCast func(string) (T, error), defaultValue T) (T, error) {
	val, exists := os.LookupEnv(variableName)
	if !exists {
		return defaultValue, nil
	}
	castVal, err := fnCast(val)
	if err != nil {
		return defaultValue, err
	}
	return castVal, nil
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

func CastInt(val string) (int, error) {
	iVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return iVal, nil
}
