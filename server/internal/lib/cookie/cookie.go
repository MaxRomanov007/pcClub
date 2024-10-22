package cookie

import (
	"net/http"
	"time"
)

func Set(w http.ResponseWriter, name string, value string, path string, ttl time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
		Secure:   true,
		HttpOnly: true,
		Path:     path,
	}
	http.SetCookie(w, cookie)
}
