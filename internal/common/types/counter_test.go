package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounterFromString(t *testing.T) {
	expected := Counter(123)
	val, err := NewCounterFromString("123")
	assert.Nil(t, err)
	assert.Equal(t, expected, val)
}

func TestNewCounterFromStringError(t *testing.T) {
	_, err := NewCounterFromString("abc")
	assert.NotNil(t, err)

	_, err = NewCounterFromString("1.234")
	assert.NotNil(t, err)

	_, err = NewCounterFromString("-1234")
	assert.NotNil(t, err)
}

func TestCounterToString(t *testing.T) {
	assert.Equal(t, "123", Counter(123).String())
}
