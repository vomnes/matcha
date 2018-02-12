package account

import "testing"

var usernameTests = []struct {
	username    string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"abcdef", true, "Only lowercase letters"},
	{"123456789", true, "Only digit character"},
	{"ABCDEF", true, "Only uppercase letters"},
	{"abcdef-123456789_ABC.DEF", true, "Lowercase, uppercase, digit, -, _, ."},
	{"v", false, "Too short"},
	{"abcdef-123456789_ABCDEF.vomnes.vomnes.vomnes.vomnes.vomnes.vomnes", false, "Too long"},
	{"vomnes->#$%^&*()_)(*&^%)", false, "Forbidden characters"},
	{" abcABC123 ", false, "Space before and after"},
}

func TestIsValidUsername(t *testing.T) {
	for _, tt := range usernameTests {
		actual := IsValidUsername(tt.username)
		if actual != tt.expected {
			t.Errorf("IsValidUsername(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.username, tt.expected, actual, tt.testContent)
		}
	}
}

var nameTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"abcdef", true, "Only lowercase letters"},
	{"ABCDEF", true, "Only uppercase letters"},
	{"abcdefABCDEF", true, "Lowercase, uppercase characters"},
	{"abcdef-ABCDEF", true, "Lowercase, uppercase and separator (HyphenMinus) characters"},
	{"abcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcd", true, "Limit max length"},
	{"", false, "Too short"},
	{"123456789", false, "Only digit character"},
	{"abcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcdefABCDEFabcde", false, "Too long - Over max length"},
	{"vomnes.", false, "Forbidden characters '.'"},
	{"vomnes_", false, "Forbidden characters '_'"},
	{"vomnes%", false, "Forbidden characters '%'"},
	{"vomnes<", false, "Forbidden characters '<'"},
	{"vomnes>", false, "Forbidden characters '>'"},
	{"vomnes=/*-+/=%^&*()", false, "Forbidden characters"},
	{" abcABC123 ", false, "Space before and after"},
}

func TestIsValidFirstLastName(t *testing.T) {
	for _, tt := range nameTests {
		actual := IsValidFirstLastName(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidFirstLastName(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var emailAddressTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"valentin.omnes@gmail.com", true, "Valid email address"},
	{"vomnes@student.42.fr", true, "Valid email address"},
	{"true@student-42-fr", true, "Valid email address"},
	{"true_true@student.fr", true, "'_' before '@'"},
	{"f@s.fr", true, "Short email address"},
	{"ç$€§/az@gmail.com", false, "Illegal characters before '@'"},
	{"false@student_42_fr", false, "Illegal characters '_'"},
	{"false@student<42.fr", false, "Illegal characters '<'"},
	{"false@student>42.fr", false, "Illegal characters '>'"},
	{"false@student@42.fr", false, "Illegal characters '@'"},
	{"false@student*42.fr", false, "Illegal characters '*'"},
	{"studentstudentstudentstudentstudentstudentstudentstudentstudentstudent" +
		"studentstudent@studentstudentstudentstudentstudentstudentstudentstudent" +
		"studentstudentstudentstudent.fr", false, "Too long email address"},
}

func TestIsValidEmailAddress(t *testing.T) {
	for _, tt := range emailAddressTests {
		actual := IsValidEmailAddress(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidEmailAddress(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}

var passwordTests = []struct {
	str         string // input
	expected    bool   // expected result
	testContent string // test details
}{
	{"abcABC123", true, "Valid password"},
	{"abcABC12", true, "Valid password 8 characters"},
	{"abcABC123abcABC123abcABC123abcABC123abcABC123abcABC" +
		"abcABC123abcABC123abcABC123abcABC123abcABC123abcA", true, "Valid password" +
		" 100 characters"},
	{"abcdefgh", false, "Only lowercase letters"},
	{"ABCDEFGH", false, "Only uppercase letters"},
	{"123456789", false, "Only digits"},
	{"abcABC1", false, "Too short"},
	{"abcABC123abcABC123abcABC123abcABC123abcABC123abcABC" +
		"abcABC123abcABC123abcABC123abcABC123abcABC123abcABC", false, "Too long"},
	{"     \t      ", false, "Only space"},
	{"", false, "Empty"},
	{"abcABC123$%^&*()_", false, "Special characters"},
}

func TestIsValidPassword(t *testing.T) {
	for _, tt := range passwordTests {
		actual := IsValidPassword(tt.str)
		if actual != tt.expected {
			t.Errorf("IsValidPassword(%s): expected %t, actual %t - Test type: \033[31m%s\033[0m", tt.str, tt.expected, actual, tt.testContent)
		}
	}
}
