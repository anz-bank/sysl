package importer

// Response is either going to be freetext or a type
type Response struct {
	Text    string
	Type    Type
	Content content
}
type content struct {
	name        string
	contentType string
}

// Endpoints consists of all the values necessary to make a sysl file
type Endpoint struct {
	Path        string
	Description string
	Params      Parameters
	Responses   []Response
}

// nolint:gochecknoglobals
var methodDisplayOrder = []string{"GET", "PUT", "POST", "DELETE", "PATCH"}
