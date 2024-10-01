package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"songs-library-go/internal/domain"
	"time"
)

func Init() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("customDate", customDateValidation)

	return validate
}

func customDateValidation(fl validator.FieldLevel) bool {
	if fl.Field().Kind() == reflect.Ptr && fl.Field().IsNil() {
		return true
	}

	var dateStr string
	if fl.Field().Kind() == reflect.Ptr {
		dateStr = fl.Field().Elem().String()
	} else {
		dateStr = fl.Field().String()
	}

	_, err := time.Parse(domain.DateFormat, dateStr)
	return err == nil
}
