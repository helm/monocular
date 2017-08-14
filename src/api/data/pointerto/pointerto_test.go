package pointerto

import (
	"testing"

	"github.com/arschles/assert"
)

func TestInt64(t *testing.T) {
	var number int64
	number = 13
	ptr := Int64(number)
	assert.Equal(t, number, *ptr, "int64 to ptr conversion")
}

func TestString(t *testing.T) {
	var str string
	str = "string"
	ptr := String(str)
	assert.Equal(t, str, *ptr, "string to ptr conversion")
}
