package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/mux"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/isutare412/imageer/internal/gateway/webv2/gen"
	"github.com/isutare412/imageer/pkg/apperr"
)

func WithOpenAPIValidator() mux.MiddlewareFunc {
	swagger, err := gen.GetSwagger()
	if err != nil {
		panic(fmt.Errorf("getting swagger spec: %w", err))
	}

	// NOTE: Manually add docs endpoints to swagger to bypass validation
	swagger.AddOperation("/docs", "GET", openapi3.NewOperation())
	swagger.AddOperation("/docs/openapi.html", "GET", openapi3.NewOperation())
	swagger.AddOperation("/docs/openapi.yaml", "GET", openapi3.NewOperation())

	// Disable schema error details in responses
	openapi3.SchemaErrorDetailsDisabled = true

	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		ErrorHandlerWithOpts:  handleOpenAPIError,
		SilenceServersWarning: true,
		DoNotValidateServers:  true,
	})
}

func handleOpenAPIError(ctx context.Context, err error, w http.ResponseWriter, r *http.Request,
	opts oapimiddleware.ErrorHandlerOpts,
) {
	var (
		summary = "Failed to validate request"
		detail  = err.Error()
	)

	var reqErr *openapi3filter.RequestError
	if errors.As(err, &reqErr) {
		summary = reqErr.Reason
		detail = reqErr.Err.Error()
	}

	gen.RespondError(w, r, apperr.NewError(apperr.DefaultCode(opts.StatusCode)).
		WithSummary("%s", summary).
		WithDetail("%s", detail))
}
