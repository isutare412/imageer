package validation

import (
	"fmt"
	"reflect"
	"strings"

	locale_us "github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	trans_en "github.com/go-playground/validator/v10/translations/en"

	"github.com/isutare412/imageer/pkg/apperr"
	"github.com/isutare412/imageer/pkg/validation/custom"
)

var globalValidator = NewValidator()

// Validate validates a struct using validator tags. This is a shortcut for
// creating a Validator and calling its Validate method.
// If you need custom validation behavior, create a Validator using
// [NewValidator].
func Validate(i any) error {
	return globalValidator.Validate(i)
}

// Validator validates structs using tags.
// Ref: https://pkg.go.dev/github.com/go-playground/validator/v10
type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewValidator() Validator {
	validate := validator.New(validator.WithRequiredStructEnabled())
	translator := ut.New(locale_us.New()).GetFallback()

	// Support english error translations
	if err := trans_en.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(fmt.Errorf("registering translator: %w", err))
	}

	// Use some tags as field name if exists
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		switch name {
		case "", "-": // try next tag
		default:
			return name
		}

		name = strings.SplitN(fld.Tag.Get("koanf"), ",", 2)[0]
		switch name {
		case "", "-": // try next tag
		default:
			return name
		}
		return ""
	})

	// Register custom tags
	if err := custom.RegisterKebabCaseTag(validate, translator); err != nil {
		panic(fmt.Errorf("registering kebabcase tag: %w", err))
	}

	return Validator{
		validate:   validate,
		translator: translator,
	}
}

func (v Validator) Validate(i any) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)
	errorMsgs := make([]string, 0, len(errs))
	for key, msg := range errs.Translate(v.translator) {
		// Remove top-level struct name from the key
		// e.g. Request.user.name -> user.name
		if i := strings.Index(key, "."); i != -1 {
			key = key[i+1:]
		}

		msg := fmt.Sprintf("%s: %s", key, msg)
		errorMsgs = append(errorMsgs, msg)
	}

	return apperr.NewError(apperr.CodeBadRequest).
		WithSummary("Validation failed").
		WithDetail("%s", strings.Join(errorMsgs, "\n"))
}
