package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mivinci/log"
	"github.com/mivinci/ttl"
)

const (
	HeaderToken = "GC-Auth-Token"
	defTokenTTL = time.Hour * 24 * 7
	defCodeTTL  = time.Minute * 10
)

var (
	ErrID              = errors.New("empty id provided")
	ErrKey             = errors.New("empty key provided")
	ErrQuery           = errors.New("incompleted get query")
	ErrParam           = errors.New("incompleted parameter")
	ErrAuthHeader      = errors.New("empty auth header")
	ErrHTTPMethod      = errors.New("bad http method")
	ErrCodeMatch       = errors.New("code is not correct")
	ErrPermDenied      = errors.New("permission denied")
	ErrAccessForbidden = errors.New("forbidden access")
	ErrMetadata        = errors.New("no metadata in request context")
)

type Auth struct {
	key string

	opt Option
}

func New(key string, opts ...Options) *Auth {
	opt := Option{
		tokenTTL:  defTokenTTL,
		codeTTL:   defCodeTTL,
		sender:    &sender{},
		generater: &generater{},
	}
	for _, o := range opts {
		o(&opt)
	}
	return &Auth{key: key, opt: opt}
}

func (a *Auth) PublicFunc(next http.HandlerFunc, excludes ...string) http.HandlerFunc {
	return a.protect(next, false, excludes...)
}

func (a *Auth) PrivateFunc(next http.HandlerFunc, excludes ...string) http.HandlerFunc {
	return a.protect(next, true, excludes...)
}

func (a *Auth) protect(next http.HandlerFunc, private bool, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !hasMethod(r, methods...) {
			var err error
			log.Debug("auth.protect: ", r.URL.Path)
			if r, err = a.OK(r, private); err != nil {
				log.Debug("auth.OK: ", err)
				a.opt.render.Error(w, err) // 401
				return
			}
		}
		next(w, r)
	}
}

func (a *Auth) TransportFunc(next http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !hasMethod(r, methods...) {
			a.transport(w, r, next)
			log.Debug("auth: transported request ", r.URL.Path)
			return
		}
		next(w, r)
	}
}

func (a *Auth) transport(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	to, code := r.PostFormValue("to"), r.PostFormValue("code")
	if to == "" || code == "" {
		a.opt.render.Error(w, ErrQuery) // 400
		return
	}

	// c, err := ttl.GetAndRemove(codeKey(to))
	// if err != nil {
	// 	log.Error("auth.Login: get code from cache error due to", err)
	// 	a.opt.render.Error(w, err)
	// 	return
	// }
	// if c != code {
	// 	a.opt.render.Error(w, ErrCodeMatch)
	// 	return
	// }

	token, err := a.opt.generater.Token(to)
	if err != nil {
		log.Error("auth.Login: generate token error due to ", err)
		a.opt.render.Error(w, err)
		return
	}

	// in case that caller fails to set secret key
	if token.Key == nil {
		token.Key = []byte(a.key)
	}

	final := token.Final()
	log.Debug("auth.transport: final token: ", final)

	r = WithContext(r, Metadata{HeaderToken: token})
	r.Header.Set(HeaderToken, final)
	w.Header().Set(HeaderToken, final)

	next(w, r)
}

// Public should be used as an HTTP middleware
func (a *Auth) Public(next http.Handler, excludes ...string) http.HandlerFunc {
	return a.PublicFunc(next.ServeHTTP, excludes...)
}

// Private should be used as an HTTP middleware
func (a *Auth) Private(next http.Handler, excludes ...string) http.HandlerFunc {
	return a.PrivateFunc(next.ServeHTTP, excludes...)
}

// Transport should be used as an HTTP middleware
func (a *Auth) Transport(next http.Handler, excludes ...string) http.HandlerFunc {
	return a.TransportFunc(next.ServeHTTP, excludes...)
}

func (a *Auth) Code(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	log.Debug("sending code to: ", to)
	if to == "" {
		a.opt.render.Error(w, ErrQuery) // 400
		return
	}

	var err error
	var code string

	if code, err = a.opt.generater.Code(); err != nil {
		log.Error("auth.Code: create code error due to", err)
		a.opt.render.Error(w, err) // 500
		return
	}

	if err = ttl.Add(codeKey(to), code, a.opt.codeTTL); err != nil {
		log.Error("auth.Code: store code error due to", err)
		a.opt.render.Error(w, err) // 500
		return
	}
	// user can never know if code is sent successfully
	go a.opt.sender.Send(to, code) // nolint: errcheck
	a.opt.render.OK(w)
}

func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.Code(w, r)
	default:
		a.opt.render.Error(w, ErrHTTPMethod)
	}
}

func (a *Auth) OK(r *http.Request, private bool) (*http.Request, error) {
	var err error
	tk := r.Header.Get(HeaderToken)
	if tk == "" {
		cookie, err := r.Cookie(HeaderToken)
		if err != nil {
			return r, err
		}
		tk = cookie.Value
		log.Debug("auth: get token from cookie")
	}

	log.Debug("receive token: ", tk)

	token := Token{Key: []byte(a.key)}
	if err = token.OK(tk); err != nil {
		return r, err
	}

	if private {
		if r, err = a.opt.filter.Private(r, token); err != nil {
			return r, err
		}
	}
	return r, nil
}

func (a Auth) CodeTTL() time.Duration  { return a.opt.codeTTL }
func (a Auth) TokenTTL() time.Duration { return a.opt.tokenTTL }

func codeKey(to string) string {
	return fmt.Sprintf("gc-code:%s", to)
}

func hasMethod(r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if method == r.Method {
			return true
		}
	}
	return false
}
