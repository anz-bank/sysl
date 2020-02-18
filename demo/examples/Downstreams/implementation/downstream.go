// downstream.go contains the implementation of the downstream app
//
//

package implementation

import (
	"context"

	"github.service.anz/sysl/syslbyexample/_examples/Downstreams/gen/downstream"
)

type ExampleService struct {
}

// GetBarList is the stubbed method of our downstream call
func (s *ExampleService) GetBarList(ctx context.Context, req *downstream.GetBarListRequest) (*downstream.Foo, error) {
	return &downstream.Foo{Content: 1234567}, nil
}
