package request

import (
	"errors"
	"net/http"
)

var (
	ErrUIDNotExists = errors.New("uid was not in request context")
)

// MustUID gets the uid from the request context and
// panics when the uid is not in the request context
func MustUID(r *http.Request) int64 {
	uid, err := UID(r)
	if err != nil {
		panic(err)
	}
	return uid
}

// UID gets the uid from the request context
func UID(r *http.Request) (int64, error) {
	uid := r.Context().Value("uid")
	if uid == nil {
		return 0, ErrUIDNotExists
	}
	return r.Context().Value("uid").(int64), nil
}
