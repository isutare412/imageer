package validation

import (
	"fmt"
	"reflect"
	"strings"

	locale_us "github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	trans_en "github.com/go-playground/validator/v10/translations/en"
)

var globalTranslator ut.Translator

func init() {
	globalTranslator = ut.New(locale_us.New()).GetFallback()
}

func NewValidate() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Support english error translations
	if err := trans_en.RegisterDefaultTranslations(v, globalTranslator); err != nil {
		panic(fmt.Errorf("registering translator: %w", err))
	}

	// Use json tag as field name if exists
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

func TranslateError(err validator.FieldError) string {
	return err.Translate(globalTranslator)
}
