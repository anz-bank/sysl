// implementation/methods.go

// methods.go contains the actual HTTP functions that control the behaviour for the endpoints defined in the sysl file
package implementation

import (
	"context"

	// the simple package is the code generated sysl in our gen directory
	downstream "github.service.anz/sysl/syslbyexample/_examples/Downstreams/gen/downstream"
	simple "github.service.anz/sysl/syslbyexample/_examples/Downstreams/gen/simple"
)

// GetStuffList refers to the first GET endpoint we specify in our sysl model
func GetStuffList(ctx context.Context, req *simple.GetStuffListRequest, client simple.GetStuffListClient) (*simple.Stuff, error) {

	return &simple.Stuff{Content: "Hello World!"}, nil
}

// GetFoobarList refers to the second GET endpoint we specify in our sysl model
func GetFoobarList(ctx context.Context, req *simple.GetFoobarListRequest, client simple.GetFoobarListClient) (*downstream.Foo, error) {
	// Call GetBarList which is contained in our downstream client
	return client.GetBarList(ctx, &downstream.GetBarListRequest{})
}
