package auth

import (
	"context"
	"net/http"
)

type Metadata map[string]interface{}

type HeaderMetadata struct{}

func FromContext(r *http.Request) (md Metadata, ok bool) {
	md, ok = r.Context().Value(HeaderMetadata{}).(Metadata)
	return
}

func WithContext(r *http.Request, md Metadata) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), HeaderMetadata{}, md))
}
