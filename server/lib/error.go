package lib

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

func PrettyError(err string) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New("Unidentified: " + err)
	}
	lastIndex := strings.LastIndex(file, "src")
	fileBytes := []byte(file)
	fileBytes = fileBytes[lastIndex:]
	file = string(fileBytes)
	return errors.New(file + ": l." + strconv.Itoa(line) + ": " + err)
}
