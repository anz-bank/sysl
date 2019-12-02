package parse

import "fmt"

type Exit struct {
	Code    int
	message string
}

func (e Exit) Error() string {
	return e.message
}

func Exitf(code int, format string, args ...interface{}) Exit {
	return Exit{code, fmt.Sprintf(format, args...)}
}
