package settings

import (
	"testing"

	"github.com/devldavydov/promytheus/internal/common/env"
	"github.com/stretchr/testify/assert"
)

func TestGetPriorityParam(t *testing.T) {
	assert.Equal(t, "valBoth", GetPriorityParam(&env.EnvPair[string]{Value: "valBoth", IsDefault: false}, "valBoth"))
	assert.Equal(t, "valBoth", GetPriorityParam(&env.EnvPair[string]{Value: "valBoth", IsDefault: true}, "valBoth"))
	assert.Equal(t, "valE", GetPriorityParam(&env.EnvPair[string]{Value: "valE", IsDefault: false}, "valF"))
	assert.Equal(t, "valF", GetPriorityParam(&env.EnvPair[string]{Value: "valE", IsDefault: true}, "valF"))
}
