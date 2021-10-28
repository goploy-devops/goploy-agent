package core

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"reflect"
	"strings"
)

// Validate use a single instance of Validate, it caches struct info
var Validate *validator.Validate

// Trans Translator
var Trans ut.Translator

// CreateValidator create a single Validator
func CreateValidator() {

	english := en.New()
	uni := ut.New(english, english)
	Trans, _ = uni.GetTranslator("english")
	Validate = validator.New()
	en_translations.RegisterDefaultTranslations(Validate, Trans)
	registerTagName()
}

func registerTagName() {
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
