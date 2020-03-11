// Package catalog takes a sysl module with attributes defined (catalogFields)
// and serves a webserver listing the applications and endpoints
// It also uses gRPC UI and Redoc in order to generate an interactive page to interact with all the endpoints
// gRPC currently uses server reflection TODO: Support gpcui directly from proto files
package catalog

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/fullstorydev/grpcui/standalone"
	"github.com/fullstorydev/grpcurl"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// APIDoc contains everything required to serve a view of a Sysl described application
type APIDoc struct {
	name    string       // App name
	spec    []byte       // JSON API Spec
	handler http.Handler // Handler to serve docs
	path    string       // URL Path where the API documentation will be served
}

func (a *APIDoc) GetHandler() http.Handler {
	return a.handler
}

type APIDocBuilder struct {
	app    *sysl.Application
	doc    *APIDoc
	log    *logrus.Logger
	grpcui bool
}

func MakeAPIDocBuilder(app *sysl.Application, log *logrus.Logger, grpcui bool) *APIDocBuilder {
	return &APIDocBuilder{
		app:    app,
		log:    log,
		grpcui: grpcui,
	}
}

//nolint:stylecheck
func (b APIDocBuilder) BuildAPIDoc() (*APIDoc, error) {
	b.doc = &APIDoc{
		name: b.app.GetName().GetPart()[0],
		path: makePath(serviceType(b.app), "/", b.app.GetName().GetPart()[0]),
	}
	err := b.generateHandlers()

	if err != nil {
		b.log.Infof("Failed to add %s:  %s\n",
			b.app.GetName().GetPart()[0],
			err)
		return nil, err
	}

	b.log.Infof("Added %s: ", b.app.GetName().GetPart()[0])
	return b.doc, nil
}

func (b APIDocBuilder) generateHandlers() error {
	var err error
	switch serviceType(b.app) {
	case GRPC:
		err = b.buildGrpcHandler()
	case REST:
		err = b.buildRestHandler(makePath(serviceType(b.app), "/", b.app.GetName().GetPart()[0]))
	default:
		b.log.Debugf("Skipping '%s' as no API type has been tagged", b.app.GetName().GetPart()[0])
	}
	return err
}

// GrpcUIHandler creates and returns a http handler for a grpcui server
func (b *APIDocBuilder) buildGrpcHandler() error {
	if !b.grpcui {
		b.log.Infoln("Skipping grpc handler")
		return nil
	}
	ctx := context.Background()
	dialURL := b.app.GetAttrs()["deploy.prod.url"].GetS()
	cc, err := buildGrpcClientConnection(ctx, dialURL)
	if err != nil {
		b.log.Errorf("Failed to create a gRPC connection to %s", dialURL)
	}
	h, err := standalone.HandlerViaReflection(ctx, cc, dialURL)
	if err != nil {
		b.log.Infoln("Failed to generate handler via reflection.")
		return err
	}

	noTrailingSlash := b.doc.path[:len(b.doc.path)-1]
	b.doc.handler = http.StripPrefix(noTrailingSlash, h)
	return nil
}

func buildGrpcClientConnection(ctx context.Context, dialURL string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	creds, err := grpcurl.ClientTransportCredentials(false, "", "", "")
	if err != nil {
		return nil, err
	}
	cc, err := grpcurl.BlockingDial(ctx, "tcp", dialURL, creds, opts...)

	// If that failed, try an insecure dial
	if err != nil {
		cc, err = grpc.DialContext(ctx, dialURL, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	}
	return cc, nil
}

func makePath(serviceMethod string, basepath string, serviceName string) string {
	// We need a trailing slash for gRPCUI to work correctly.
	// TODO: Fix this and remove this hacky workaround
	trailingSlash := ""
	if serviceMethod == GRPC {
		trailingSlash = "/"
	}

	return path.Join(basepath, strings.ToLower(serviceMethod), serviceName) + trailingSlash
}
