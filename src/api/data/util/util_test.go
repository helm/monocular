package util

import (
	"testing"

	"github.com/arschles/assert"
)

func TestInt64ToPtr(t *testing.T) {
	var number int64
	number = 13
	ptr := Int64ToPtr(number)
	assert.Equal(t, number, *ptr, "int64 to ptr conversion")
}

func TestStrToPtr(t *testing.T) {
	var str string
	str = "string"
	ptr := StrToPtr(str)
	assert.Equal(t, str, *ptr, "string to ptr conversion")
}
