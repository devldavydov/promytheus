package settings

func GetPriorityParam[T comparable](envValue, flagValue, defaultValue T) T {
	if envValue == flagValue {
		return envValue
	}

	if envValue != defaultValue {
		return envValue
	}

	return flagValue
}
