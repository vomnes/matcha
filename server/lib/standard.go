package lib

type key int

const (
	// Database key is used as value in the context
	Database = key(1)
	// UserID key is used as value in the context
	UserID = key(2)
)

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

func strsub(input string, start int, end int) string {
	var output string
	for i := start; i < start+end; i++ {
		output += string(input[i])
	}
	return output
}
