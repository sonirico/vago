package str

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return s == ""
}

// IsSet checks if a string is set (not empty)
func IsSet(s string) bool {
	return s != ""
}
