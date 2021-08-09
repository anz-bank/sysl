package syslutil

import "github.com/anz-bank/sysl/pkg/sysl"

// nolint:revive,stylecheck
var (
	Method_GET    = sysl.Endpoint_RestParams_GET.String()
	Method_PUT    = sysl.Endpoint_RestParams_PUT.String()
	Method_POST   = sysl.Endpoint_RestParams_POST.String()
	Method_DELETE = sysl.Endpoint_RestParams_DELETE.String()
	Method_PATCH  = sysl.Endpoint_RestParams_PATCH.String()
)
