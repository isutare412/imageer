package web

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

func openAPIValidator() echo.MiddlewareFunc {
	swagger, err := GetSwagger()
	if err != nil {
		panic(fmt.Errorf("getting swagger spec: %w", err))
	}

	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
		Skipper:               skipOpenAPIValidation,
		SilenceServersWarning: true,
	})
}

func skipOpenAPIValidation(c echo.Context) bool {
	return !strings.HasPrefix(c.Request().URL.Path, "/api")
}
