package handler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/mivinci/s/auth"
)

type Transport struct {
	Auth *auth.Auth
}

func (t *Transport) Login(w http.ResponseWriter, r *http.Request) {
	md, ok := auth.FromContext(r)
	if !ok {
		Error(w, auth.ErrMetadata)
		return
	}

	token, ok := md[auth.HeaderToken].(auth.Token)
	if !ok {
		Error(w, ErrAuthHeader)
		return
	}

	cookie := http.Cookie{
		Name:   auth.HeaderToken,
		Value:  token.Final(),
		Path:   "/",
		MaxAge: int(t.Auth.TokenTTL() / time.Second),
	}

	http.SetCookie(w, &cookie)

	ref, _ := url.ParseRequestURI(r.Referer())
	redirect := ref.Query().Get("redirect")
	if redirect == "" {
		redirect = r.PostFormValue("redirect")
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}

func (t *Transport) Code(w http.ResponseWriter, r *http.Request) {
	t.Auth.Code(w, r)
}

func (t *Transport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t.Code(w, r)
	case "POST":
		t.Login(w, r)
	default:
		Error(w, ErrHTTPMethod)
	}
}
