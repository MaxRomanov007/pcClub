package validator

import (
	"fmt"
	"github.com/go-playground/validator"
	"reflect"
	reflect2 "server/internal/lib/reflect"
	"strings"
)

func ValidationError[T any](errs validator.ValidationErrors) string {
	var errMsgs []string

	t := reflect2.TypeOf[T]()

	for _, err := range errs {
		var field string
		field = fieldName(err, t)

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is required", field))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is over maximum", field))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is lover minimum", field))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("the field %s is not valid", field))
		}
	}
	return strings.Join(errMsgs, "; ")
}

func fieldName(err validator.FieldError, t reflect.Type) string {
	field := trimLeft(err.StructNamespace())
	if trimLeft(t.String()) != trimRight(err.StructNamespace()) {
		return field
	}

	f, ok := t.FieldByName(field)
	if !ok {
		return field
	}

	for _, tagName := range []string{"json", "get"} {
		name, ok := f.Tag.Lookup(tagName)
		if !ok {
			continue
		}
		return name
	}
	return field
}

func trimLeft(str string) string {
	index := strings.Index(str, ".")
	if index == -1 {
		return str
	}
	return str[index+1:]
}

func trimRight(str string) string {
	index := strings.Index(str, ".")
	if index == -1 {
		return str
	}
	return str[:index]
}
