package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/isutare412/imageer/pkg/apperr"
)

// StructValidator is a custom validator for struct validation.
type StructValidator struct {
	validate *validator.Validate
}

func NewStructValidator(v *validator.Validate) StructValidator {
	return StructValidator{
		validate: v,
	}
}

func (v StructValidator) Validate(i any) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	errorMsgs := make([]string, 0, len(errs))
	for _, e := range errs {
		errorMsgs = append(errorMsgs, TranslateError(e))
	}

	return apperr.NewError(apperr.CodeBadRequest).
		WithSummary("Validation failed").
		WithDetail("%s", strings.Join(errorMsgs, "\n"))
}
