package main

const (
	// ParserSuccess is returned by parser when it was able to parse input correctly
	ParserSuccess = 0
	// ImportError is returned by parser when its unable to load input modules
	ImportError = 1
	// ParseError is returned by parser when one of the input files has syntax errors
	ParseError = 2
)
