// main.go
package main

import (
	"context"
	"github.service.anz/sysl/syslbyexample/_examples/LiveDownstream/implementation"
)

func main() {
	// Now the LoadServices function is called to start our server
	implementation.LoadServices(context.Background())
}
