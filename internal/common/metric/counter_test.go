package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounterFromString(t *testing.T) {
	expected := Counter(123)
	val, err := NewCounterFromString("123")
	assert.NoError(t, err)
	assert.Equal(t, expected, val)
}

func TestNewCounterFromStringError(t *testing.T) {
	_, err := NewCounterFromString("abc")
	assert.Error(t, err)

	_, err = NewCounterFromString("1.234")
	assert.Error(t, err)

	_, err = NewCounterFromString("-1234")
	assert.Error(t, err)
}

func TestNewCounterFromIntPError(t *testing.T) {
	_, err := NewCounterFromIntP(nil)
	assert.Error(t, err)
}

func TestCounterToIntP(t *testing.T) {
	var exp int64 = 123
	assert.Equal(t, exp, *Counter(123).IntP())
}

func TestCounterToString(t *testing.T) {
	assert.Equal(t, "123", Counter(123).String())
}
