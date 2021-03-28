package importer

import "github.com/anz-bank/sysl/pkg/syslutil"

// Response is either going to be freetext or a type
type Response struct {
	Text string
	Type Type
}

type Endpoint struct {
	Path        string
	Description string

	Params Parameters

	Responses []Response
}

// nolint:gochecknoglobals
var methodDisplayOrder = []string{
	syslutil.Method_GET,
	syslutil.Method_PUT,
	syslutil.Method_POST,
	syslutil.Method_DELETE,
	syslutil.Method_PATCH,
}
