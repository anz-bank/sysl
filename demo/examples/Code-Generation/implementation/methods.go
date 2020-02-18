// implementation/methods.go

// methods.go contains the actual HTTP functions that control the behaviour for the endpoints defined in the sysl file
package implementation

import (
	"context"

	// the simple package is the code generated sysl in our gen directory
	simple "github.service.anz/sysl/syslbyexample/_examples/Code-Generation/gen/simple"
)

// GetTestList refers to the first GET endpoint we specify in our sysl model
func GetTestList(ctx context.Context, req *simple.GetTestListRequest, client simple.GetTestListClient) (*simple.Stuff, error) {
	return &simple.Stuff{Content: "9q834noxfo348t908qdjr9w8no4rmq34r"}, nil
}

// GetFoobarList refers to the second GET endpoint we specify in our sysl model
func GetFoobarList(ctx context.Context, req *simple.GetFoobarListRequest, client simple.GetFoobarListClient) (*simple.Foo, error) {
	return &simple.Foo{Content: 123123}, nil
}
