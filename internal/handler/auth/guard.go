package auth

import (
	"net/http"
)

func Guard(next http.HandlerFunc) http.HandlerFunc {

	return next
}
