// Package util contains small utilities used in other packages
package util

// Int64ToPtr converts an int64 to an *int64
func Int64ToPtr(n int64) *int64 {
	return &n
}

// StrToPtr converts a string to a *string
func StrToPtr(s string) *string {
	return &s
}
