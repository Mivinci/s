package auth

import (
	"net/http"
)

type Render interface {
	JSON(http.ResponseWriter, interface{}, error)
	Error(http.ResponseWriter, error)
	OK(http.ResponseWriter)
}
