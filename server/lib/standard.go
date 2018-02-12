package lib

// StringInArray take a string and a array of string as parameter
// Return true if the string is in the array of string else false
func StringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
