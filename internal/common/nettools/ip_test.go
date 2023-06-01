package nettools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostIP(t *testing.T) {
	res, err := GetHostIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
