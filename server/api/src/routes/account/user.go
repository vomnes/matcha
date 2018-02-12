package account

import "regexp"

const (
	// UsernameMinLength corresponds to the minimum character length of the username
	UsernameMinLength = 6
	// UsernameMaxLength corresponds to the maximum character length of the username
	UsernameMaxLength = 24
	// EmailAddressMaxLength corresponds to the maximum character length of the email address
	EmailAddressMaxLength = 255
)

// IsValidUsername check if the string parameter is a valid username
// Check length maximum and minimum and authorized characters (a-zA-Z0-9.-_)
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

// IsValidFirstLastName check if the string parameter is a valid lastname or firstname
// Check length maximum and minimum and authorized characters (a-zA-Z -)
// Return a boolean
func IsValidFirstLastName(s string) bool {
	if len(s) < UsernameMinLength || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[a-zA-Z\-]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsValidEmailAddress check if the string parameter is a valid email address
// Return a boolean
func IsValidEmailAddress(s string) bool {
	reEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !reEmail.MatchString(s) {
		return false
	}
	return true
}

// IsValidPassword check if the string parameter is a valid password
// Return a boolean
func IsValidPassword(s string) bool {
	reEmail := regexp.MustCompile("^*\\d*[a-z]*[A-Z][a-zA-Z0-9]{8, 254}$")
	if !reEmail.MatchString(s) {
		return false
	}
	return true
}
