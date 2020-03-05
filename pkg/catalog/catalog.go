// Package catalog takes a sysl module with attributes defined (catalogFields)
// and serves a webserver listing the applications and endpoints
// It also uses gRPC UI and Redoc in order to generate an interactive page to interact with all the endpoints
// gRPC currently uses server reflection TODO: Support gpcui directly from proto files
package catalog

import (
	"net/http"
	"path"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

// APIDoc contains everything required to serve a view of a Sysl described application
type APIDoc struct {
	name    string       // App name
	spec    []byte       // JSON API Spec
	handler http.Handler // Handler to serve docs
	path    string
}

func (a *APIDoc) GetHandler() http.Handler {
	return a.handler
}

type APIDocBuilder struct {
	app *sysl.Application
	doc *APIDoc
	log *logrus.Logger
}

func MakeAPIDocBuilder(app *sysl.Application, log *logrus.Logger) *APIDocBuilder {
	return &APIDocBuilder{
		app: app,
		log: log,
	}
}

//nolint:stylecheck
func (b APIDocBuilder) BuildAPIDoc() (*APIDoc, error) {
	b.doc = &APIDoc{
		name: b.app.GetName().GetPart()[0],
		path: makePath(serviceType(b.app), "/", b.app.GetName().GetPart()[0]),
	}
	doc, err := b.generateHandlers()

	if err != nil {
		b.log.Infof("Failed to add %s:  %s\n",
			b.app.GetName().GetPart()[0],
			err)
		return nil, err
	}

	b.log.Infof("Added %s: ", b.app.GetName().GetPart()[0])
	return doc, nil
}

func (b APIDocBuilder) generateHandlers() (*APIDoc, error) {
	var err error
	switch serviceType(b.app) {
	case GRPC:
		err = b.buildGrpcHandler()
	case REST:
		err = b.buildRestHandler(makePath(serviceType(b.app), "/", b.app.GetName().GetPart()[0]))
	default:
		b.log.Debugf("Skipping '%s' as no API type has been tagged", b.app.GetName().GetPart()[0])
	}
	return b.doc, err
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
