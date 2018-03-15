package lib

import (
	"regexp"
	"strconv"
	"strings"
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

// IsValidDate check if the string parameter is a valid date dd/mm/yyyy
// Return a boolean and status error
func IsValidDate(s string) (bool, error) {
	if len(s) > len("mm/dd/yyyy") {
		return false, nil
	}
	if !regexp.MustCompile(`^[0-9\/]+$`).MatchString(s) {
		return false, nil
	}
	if strings.Count(s, "/") != 2 {
		return false, nil
	}
	var day, month, year string
	idx := strings.Index(s, "/")
	if idx != -1 {
		day = s[:idx]
	}
	lastIdx := strings.LastIndex(s, "/")
	if lastIdx != -1 {
		year = s[lastIdx+1:]
	}
	month = Strsub(s, idx+1, 2)
	var monthLimit = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return false, err
	}
	if monthInt < 1 || monthInt > 12 {
		return false, nil
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return false, err
	}
	if dayInt < 1 || dayInt > monthLimit[monthInt-1] {
		return false, nil
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return false, err
	}
	if yearInt <= 999 || yearInt > 9999 {
		return false, nil
	}
	return true, nil
}

// IsValidTag check if the string parameter is a valid tag
// Check length maximum and minimum and authorized characters (a-z0-9-_)
// Return a boolean
func IsValidTag(s string) bool {
	if len(s) < 1 || len(s) > UsernameMaxLength {
		return false
	}
	if !regexp.MustCompile(`^[a-z0-9\-_]+$`).MatchString(s) {
		return false
	}
	return true
}
