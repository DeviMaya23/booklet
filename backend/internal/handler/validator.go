package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type echoValidator struct {
	v *validator.Validate
}

// NewEchoValidator returns an echo.Validator backed by go-playground/validator.
// Field names in error payloads use json struct tags instead of Go field names.
func NewEchoValidator() echo.Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		p := field.Interface().(Patch[string])
		if !p.Set || p.Value == nil {
			return ""
		}
		return *p.Value
	}, Patch[string]{})
	return &echoValidator{v: v}
}

func (ev *echoValidator) Validate(i any) error {
	return ev.v.Struct(i)
}

type validationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type validationErrBody struct {
	Errors []validationFieldError `json:"errors"`
}

func validationErrResponse(err error) validationErrBody {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return validationErrBody{Errors: []validationFieldError{}}
	}
	errs := make([]validationFieldError, len(ve))
	for i, fe := range ve {
		errs[i] = validationFieldError{
			Field:   fe.Field(),
			Message: tagMessage(fe),
		}
	}
	return validationErrBody{Errors: errs}
}

func tagMessage(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must not be empty", field)
	case "oneof":
		allowed := strings.ReplaceAll(fe.Param(), " ", ", ")
		return fmt.Sprintf("%s must be one of: %s", field, allowed)
	case "uuid4", "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	default:
		return fmt.Sprintf("%s failed validation", field)
	}
}
