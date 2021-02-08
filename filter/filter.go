package filter

import (
	"net/http"
	"strconv"

	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/handler"
	"github.com/mivinci/shortid"
)

var _ auth.Filter = Filter{}

type Filter struct {
	User *handler.User
}

func (f Filter) Private(r *http.Request, token auth.Token) (*http.Request, error) {
	uid := r.URL.Query().Get("uid")
	if uid == "" {
		uid = r.PostFormValue("uid")
	}

	id, err := strconv.Atoi(uid)
	if err != nil {
		return r, auth.ErrQuery
	}

	perm, ok := token.Meta["perm"]
	if !ok {
		return r, auth.ErrTokenIncompleted
	}

	if perm.(float64) > 1 {
		log.Debug("filter.Private: pass with superior permission ", perm)
		return withContext(r, id)
	}

	if shortid.String(id) != token.ID {
		return r, auth.ErrPermDenied
	}
	log.Debug("filter.Private: pass with matched uid: ", uid)
	return withContext(r, id)
}

func withContext(r *http.Request, id int) (*http.Request, error) {
	r = auth.WithContext(r, auth.Metadata{"uid": id})
	return r, nil
}
