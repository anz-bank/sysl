package importer

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
var methodDisplayOrder = []string{"GET", "PUT", "POST", "DELETE", "PATCH"}
