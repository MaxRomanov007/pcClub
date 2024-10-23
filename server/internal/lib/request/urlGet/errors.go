package urlGet

type Error struct {
	Message string
	Code    string
}

func (e *Error) Error() string {
	return e.Message
}

const (
	ErrTypeNotStructCode    = "TypeNotStruct"
	ErrConvertFailedCode    = "ConvertFailed"
	ErrTagConvertFailedCode = "TagConvertFailed"
	ErrUnsupportedTypeCode  = "UnsupportedType"
	ErrPrivateFieldCode     = "PrivateField"
	ErrNameFieldEmptyCode   = "EmptyFieldName"
)

var (
	ErrTypeNotStruct = &Error{
		Code:    ErrTypeNotStructCode,
		Message: "input object isn't a struct",
	}
	ErrTagConvertFailed = &Error{
		Code:    ErrConvertFailedCode,
		Message: "tag convert failed",
	}
	ErrConvertFailed = &Error{
		Code:    ErrConvertFailedCode,
		Message: "convert failed",
	}
	ErrUnsupportedType = &Error{
		Code:    ErrUnsupportedTypeCode,
		Message: "unsupported type",
	}
	ErrPrivateField = &Error{
		Code:    ErrPrivateFieldCode,
		Message: "field to set is private",
	}
	ErrNameFieldEmpty = &Error{
		Code:    ErrNameFieldEmptyCode,
		Message: "field name is empty",
	}
)
