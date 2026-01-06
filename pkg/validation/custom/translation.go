package custom

import (
	"log/slog"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func defaultRegisterFunc(tag string, translation string, override bool,
) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		if err := ut.Add(tag, translation, override); err != nil {
			return err
		}
		return nil
	}
}

func defaultTranslateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		slog.Error("Failed to translate validation error", "error", err, "fieldError", fe)
		return fe.Error()
	}

	return t
}
