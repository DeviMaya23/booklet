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
			Message: tagMessage(fe.Field(), fe.Tag()),
		}
	}
	return validationErrBody{Errors: errs}
}

func tagMessage(field, tag string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must not be empty", field)
	default:
		return fmt.Sprintf("%s failed validation", field)
	}
}
