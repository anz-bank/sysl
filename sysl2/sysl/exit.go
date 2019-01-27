package main

import "fmt"

type exit struct {
	code    int
	message string
}

func exitf(code int, format string, args ...interface{}) exit {
	return exit{code, fmt.Sprintf(format, args...)}
}
