package custom

import (
	"fmt"
	"reflect"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	kebabCaseTag     = "kebabcase"
	kebabCasePattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
)

func RegisterKebabCaseTag(v *validator.Validate, ut ut.Translator) error {
	if err := v.RegisterValidation(kebabCaseTag, validateKebabCase); err != nil {
		return fmt.Errorf("registering validation: %w", err)
	}

	msg := fmt.Sprintf("{0} must match the regex pattern %s", kebabCasePattern.String())
	if err := v.RegisterTranslation(kebabCaseTag, ut,
		defaultRegisterFunc(kebabCaseTag, msg, false), defaultTranslateFunc); err != nil {
		return fmt.Errorf("registering translation: %w", err)
	}

	return nil
}

func validateKebabCase(fl validator.FieldLevel) bool {
	kind := fl.Field().Kind()
	if kind == reflect.Pointer {
		kind = fl.Field().Type().Elem().Kind()
	}

	if kind != reflect.String {
		return false
	}
	return kebabCasePattern.MatchString(fl.Field().String())
}
