package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPriorityParam(t *testing.T) {
	assert.Equal(t, "valBoth", GetPriorityParam("valBoth", "valBoth", "default"))
	assert.Equal(t, "default", GetPriorityParam("default", "default", "default"))
	assert.Equal(t, "valE", GetPriorityParam("valE", "default", "default"))
	assert.Equal(t, "valF", GetPriorityParam("default", "valF", "default"))
	assert.Equal(t, "valE", GetPriorityParam("valE", "valF", "default"))
}
