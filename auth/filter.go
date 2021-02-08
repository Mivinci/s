package auth

import "net/http"

type Filter interface {
	Private(*http.Request, Token) (*http.Request, error)
}
