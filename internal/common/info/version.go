package info

import "fmt"

func FormatVersion(version, date, commit string) string {
	_s := func(s string) string {
		if s == "" {
			return "N/A"
		}
		return s
	}

	return fmt.Sprintf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		_s(version), _s(date), _s(commit),
	)
}
