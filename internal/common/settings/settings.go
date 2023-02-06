package settings

import "github.com/devldavydov/promytheus/internal/common/env"

func GetPriorityParam[T comparable](envValue *env.EnvPair[T], flagValue T) T {
	if envValue.Value == flagValue {
		return envValue.Value
	}

	if !envValue.IsDefault {
		return envValue.Value
	}

	return flagValue
}
