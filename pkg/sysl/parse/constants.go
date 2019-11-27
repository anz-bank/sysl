package parse

const (
	// ImportError is returned by parser when its unable to load input modules
	ImportError = 1
	// ParseError is returned by parser when one of the input files has syntax errors
	ParseError = 2
)

const syslExt = ".sysl"
