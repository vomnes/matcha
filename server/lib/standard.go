package lib

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

type key int

const (
	// Database key is used as value in order to store database in the context
	Database key = iota
	// UserID key is used as value in order to store userId from JSON Web Token in the context
	UserID
	// Username key is used as value in order to store username from JSON Web Token in the context
	Username
	// Redis key is used as value in order to store redis client in the context
	Redis
	// UUID key is used as value in order to store the UUID from JSON Web Token
	// in the context, used for logout
	UUID
	// MailJet key is used as value in order to store MailJet client in the context
	MailJet
)

var (
	// PostgreSQLName is the current PostgreSQL database name
	PostgreSQLName = os.Getenv("DB_NAME")
	// PostgreSQLNameTests is the PostgreSQL database name for tests
	PostgreSQLNameTests = os.Getenv("DB_NAME_TEST")
	// RedisDBNum is the Redis database int used in the database selection
	RedisDBNum = 0
	// RedisDBNumTests is the Redis tests database int used in the database selection
	RedisDBNumTests = 1
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

func Strsub(input string, start int, end int) string {
	var output string
	for i := start; i < start+end; i++ {
		output += string(input[i])
	}
	return output
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789_-"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GetRandomString create a random string with a length of n characters
// with the characters include in letterBytes
func GetRandomString(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// UniqueTimeToken generates a unique base64 token using a key and time.Now() as string
func UniqueTimeToken(key string) string {
	now := time.Now()
	data := []byte(key + "&" + now.String())
	return base64.StdEncoding.EncodeToString(data)
}

func SWAPStrings(str1, str2 *string) {
	*str1, *str2 = *str2, *str1
}

func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func ExtractBase64Struct(base64 string, data interface{}) error {
	byteData, err := Base64Decode(base64)
	if err != nil {
		log.Println(PrettyError("[Base64] Failed to decode search parameters in header " + err.Error()))
		return err
	}
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		log.Println(PrettyError("[Unmarshal] Failed to unmarshal search parameters in header " + err.Error()))
		return err
	}
	return nil
}
