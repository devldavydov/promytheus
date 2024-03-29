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

func TestNewGaugeFromFloatPError(t *testing.T) {
	_, err := NewGaugeFromFloatP(nil)
	assert.Error(t, err)
}

func TestGaugeToFloatP(t *testing.T) {
	exp := 123.0
	assert.Equal(t, exp, *Gauge(123.0).FloatP())
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

func TestGaugeTypeName(t *testing.T) {
	assert.Equal(t, GaugeTypeName, Gauge(123.123).TypeName())
}

func TestGaugeHmac(t *testing.T) {
	assert.Equal(t,
		"2f7e24d0b4f8ff42c7ca771b44e80daaf0460ed08cbccb7228b821a9b5204934",
		Gauge(123.123).Hmac("foo", "bar"))
}
