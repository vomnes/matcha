package account

import "regexp"

const (
	// UsernameMinLength is a in corresponding to the minimum character length of the username
	UsernameMinLength = 6
	// UsernameMaxLength is a in corresponding to the maximum character length of the username
	UsernameMaxLength = 24
)

// IsValidUsername check if the string parameter is a valid username
// Check length maximum and minimum and authorized characters (a-z)
// Return a boolean
func IsValidUsername(s string) bool {
	if len(s) < UsernameMinLength || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`).MatchString(s) {
		return false
	}
	return true
}
