package filter

import (
	"net/http"

	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/metadata"
)

var _ auth.Filter = Filter{}

type Filter struct{}

func (f Filter) Private(r *http.Request, token auth.Token) (*http.Request, error) {
	perm, ok := token.Meta["perm"]
	if !ok {
		return r, auth.ErrTokenPerm
	}

	uid, ok := token.Meta["uid"]
	if !ok {
		return r, auth.ErrTokenUid
	}

	id, ok := metadata.Int(uid)
	if !ok {
		return r, auth.ErrTokenDecode
	}

	pm, ok := metadata.Int(perm)
	if !ok {
		return r, auth.ErrTokenDecode
	}

	if pm > 1 {
		log.Debug("filter.Private: pass with superior permission ", perm)
		return WithContext(r, id), nil
	}

	log.Debug("filter.Private: pass with matched uid: ", uid)
	return WithContext(r, id), nil
}

// WithContext injects uid to the flow of request
func WithContext(r *http.Request, uid int) *http.Request {
	return auth.WithContext(r, auth.Metadata{"uid": uid})
}

// FromConext extracts uid from token
func FromConext(r *http.Request) (uid int, ok bool) {
	md, ok := auth.FromContext(r)
	if !ok {
		return
	}
	id, ok := md["uid"]
	if !ok {
		return
	}
	uid, ok = id.(int)
	return
}
