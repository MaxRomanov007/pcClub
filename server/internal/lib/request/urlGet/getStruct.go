package urlGet

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func Decode(r *http.Request, s interface{}) error {
	v := reflect.ValueOf(s).Elem()

	if v.Kind() != reflect.Struct {
		return ErrTypeNotStruct
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		tag, ok := fieldType.Tag.Lookup("get")
		if !ok {
			continue
		}

		tag, isUrlParam, err := separateTag(tag)
		if err != nil {
			return err
		}
		if tag == "" {
			return ErrNameFieldEmpty
		}

		if !(field.IsValid() && field.CanSet()) {
			return fmt.Errorf("failed set field %s: %w", fieldType.Name, ErrPrivateField)
		}

		var fieldValue string
		if isUrlParam {
			fieldValue = chi.URLParam(r, tag)
		} else {
			fieldValue = r.URL.Query().Get(tag)
		}

		if fieldValue == "" {
			fieldValue = fmt.Sprintf("%v", reflect.ValueOf(reflect.Zero(fieldType.Type)).Interface())
		}

		if err := setField(field, fieldType, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func separateTag(tag string) (string, bool, error) {
	params := strings.Split(tag, ",")

	if len(params) == 1 {
		return strings.TrimSpace(params[0]), false, nil
	}
	if fieldName, ok := strings.CutPrefix(strings.TrimSpace(params[1]), "name="); ok {
		isUrlParam, err := strconv.ParseBool(strings.TrimPrefix(strings.TrimSpace(params[0]), "url="))
		if err != nil {
			return fieldName, isUrlParam, fmt.Errorf("%w: failed to convert url tag to bool", ErrTagConvertFailed)
		}
		return fieldName, isUrlParam, nil
	}
	if isUrlParamStr, ok := strings.CutPrefix(strings.TrimSpace(params[0]), "url="); ok {
		isUrlParam, err := strconv.ParseBool(strings.TrimSpace(isUrlParamStr))
		if err != nil {
			return strings.TrimPrefix(strings.TrimSpace(params[1]), "name="),
				isUrlParam,
				fmt.Errorf("%w: failed to convert url tag to bool", ErrTagConvertFailed)
		}
		return strings.TrimPrefix(strings.TrimSpace(params[1]), "name="), isUrlParam, nil
	}
	isUrlParam, err := strconv.ParseBool(strings.TrimSpace(params[1]))
	if err != nil {
		return strings.TrimSpace(params[0]),
			isUrlParam,
			fmt.Errorf("%w: failed to convert url tag to bool", ErrTagConvertFailed)
	}
	return strings.TrimSpace(params[0]), isUrlParam, nil
}

func setField(field reflect.Value, fieldType reflect.StructField, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		res, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to int8", ErrConvertFailed, fieldType.Name)
		}
		field.SetInt(int64(res))
	case reflect.Int8:
		res, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to int8", ErrConvertFailed, fieldType.Name)
		}
		field.SetInt(res)
	case reflect.Int16:
		res, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to int16", ErrConvertFailed, fieldType.Name)
		}
		field.SetInt(res)
	case reflect.Int32:
		res, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to int32", ErrConvertFailed, fieldType.Name)
		}
		field.SetInt(res)
	case reflect.Int64:
		res, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to int64", ErrConvertFailed, fieldType.Name)
		}
		field.SetInt(res)
	case reflect.Float32:
		res, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to float32", ErrConvertFailed, fieldType.Name)
		}
		field.SetFloat(res)
	case reflect.Float64:
		res, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%w: failed parse field %s to float32", ErrConvertFailed, fieldType.Name)
		}
		field.SetFloat(res)
	default:
		return fmt.Errorf("%w: unsupported field type %s", ErrUnsupportedType, fieldType.Type.Name())
	}
	return nil
}
