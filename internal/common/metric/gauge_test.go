package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGaugeFromString(t *testing.T) {
	expected := Gauge(1.23456)
	val, err := NewGaugeFromString("1.23456")
	assert.NoError(t, err)
	assert.Equal(t, expected, val)

	val, err = NewGaugeFromString("1.23456000")
	assert.NoError(t, err)
	assert.Equal(t, expected, val)
}

func TestNewGaugeFromStringErr(t *testing.T) {
	_, err := NewGaugeFromString("abc")
	assert.Error(t, err)
}

func TestGaugeToString(t *testing.T) {
	assert.Equal(t, "1.230", Gauge(1.23).String())
	assert.Equal(t, "1.235", Gauge(1.23456).String())
	assert.Equal(t, "1.000", Gauge(1).String())
}
