package pointerto

// Int64 converts an int64 to an *int64
func Int64(n int64) *int64 {
	return &n
}

// String converts a string to a *string
func String(s string) *string {
	return &s
}
