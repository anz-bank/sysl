package catalog

import (
	"context"
	"net/http"

	"github.com/fullstorydev/grpcui/standalone"
	"github.com/fullstorydev/grpcurl"
	"google.golang.org/grpc"
)

// GrpcUIHandler creates and returns a http handler for a grpcui server
func (a *APIDocBuilder) buildGrpcHandler() error {
	ctx := context.Background()
	var opts []grpc.DialOption
	var dialURL = a.app.GetAttrs()["deploy.prod.url"].GetS()
	creds, err := grpcurl.ClientTransportCredentials(false, "", "", "")
	if err != nil {
		return err
	}
	cc, err := grpcurl.BlockingDial(ctx, "tcp", dialURL, creds, opts...)
	// If that failed, try an insecure dial
	if err != nil {
		cc, err = grpc.DialContext(ctx, dialURL, grpc.WithInsecure())
		if err != nil {
			return err
		}
	}
	h, err := standalone.HandlerViaReflection(ctx, cc, dialURL)
	if err != nil {
		return err
	}

	noTrailingSlash := a.doc.path[:len(a.doc.path)-1]
	a.doc.handler = http.StripPrefix(noTrailingSlash, h)
	return nil
}
