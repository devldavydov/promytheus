// Package env provides functions to work with environment variables.
package env

import (
	"os"
	"strconv"
	"time"
)

// GetVariable retrieves variable from environment.
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

// CastString returns environment string value as string.
func CastString(val string) (string, error) {
	return val, nil
}

// CastString tries to cast environment string value to bool or error.
func CastBool(val string) (bool, error) {
	bVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}
	return bVal, nil
}

// CastDuration tries to cast environment string value to Duration or error.
func CastDuration(val string) (time.Duration, error) {
	dVal, err := time.ParseDuration(val)
	if err != nil {
		return 0, err
	}
	return dVal, nil
}

// CastDuration tries to cast environment string value to int or error.
func CastInt(val string) (int, error) {
	iVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return iVal, nil
}
