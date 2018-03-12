package lib

import (
	"regexp"
	"unicode"
)

const (
	// UsernameMinLength corresponds to the minimum character length of the username
	UsernameMinLength = 6
	// UsernameMaxLength corresponds to the maximum character length of the username
	UsernameMaxLength = 64
	// EmailAddressMaxLength corresponds to the maximum character length of the email address
	EmailAddressMaxLength = 254
	// PasswordMinLength corresponds to the minimum character length of a password
	PasswordMinLength = 8
	// PasswordMaxLength corresponds to the maximum character length of a password
	PasswordMaxLength = 100
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
	if len(s) < 1 || len(s) > UsernameMaxLength {
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
// A valid password must contains only uppercase and lowercase characters and digits,
// have a length between PasswordMinLength and PasswordMaxLength
// Return a boolean
func IsValidPassword(s string) bool {
	var digit, upper, lower bool
	char := 0
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			if !digit {
				digit = true
			}
		case unicode.IsUpper(s):
			if !upper {
				upper = true
			}
		case unicode.IsLower(s):
			if !lower {
				lower = true
			}
		default:
			return false
		}
		char++
	}
	if char < PasswordMinLength || char > PasswordMaxLength ||
		!digit || !lower || !upper {
		return false
	}
	return true
}

// IsValidText check if the string parameter is a text
// Check length maximum and authorized characters (a-zA-Z0-9 .,?!&-_*-+@#$%;)
// Return a boolean
func IsValidText(s string, lengthMax int) bool {
	if len(s) > lengthMax {
		return false
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9\ \.\,\?\!\&\-\_\*\-\+\@\#\$\%\;]+$`).MatchString(s) {
		return false
	}
	return true
}

// IsOnlyLowercaseLetters check if the string parameter has only lowercase letters
// Check length maximum and authorized characters (a-z)
// Return a boolean
func IsOnlyLowercaseLetters(s string, lengthMax int) bool {
	if len(s) > lengthMax {
		return false
	}
	if !regexp.MustCompile(`^[a-z]+$`).MatchString(s) {
		return false
	}
	return true
}

func IsValidDate(s string) bool {
	// if !regexp.MustCompile(`^(?:(?:31(\/)(?:0?[13578]|1[02]))\1|(?:(?:29|30)(\/)(?:0?[1,3-9]|1[0-2])\2))(?:(?:1[6-9]|[2-9]\d)?\d{2})$|^(?:29(\/)0?2\3(?:(?:(?:1[6-9]|[2-9]\d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1\d|2[0-8])(\/)(?:(?:0?[1-9])|(?:1[0-2]))\4(?:(?:1[6-9]|[2-9]\d)?\d{2})$`).MatchString(s) {
	// 	return false
	// }
	return true
}
