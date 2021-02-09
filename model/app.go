package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mivinci/s/auth"
	"github.com/mivinci/shortid"
)

const (
	TypeWeb int8 = iota
	TypeIOS
	TypeAndroid

	StateOK int8 = iota
	StateTerm
)

type App struct {
	ID    int       `json:"id" storm:"increment"`
	Uid   int       `json:"uid" storm:"index"`
	Type  int8      `json:"type"`
	State int8      `json:"state"`
	Name  string    `json:"name" storm:"index"`
	Site  string    `json:"site"`
	Code  string    `json:"code"` // Code is used to refresh the key of app
	Ctime time.Time `json:"ctime" storm:"index"`
	// Fields below need not be saved to db
	Key string `json:"key,omitempty"` // Key is the secret key of App
}

func (a App) Checksum(ts string) string {
	h := sha256.New()
	// encrypt secret key
	h.Write([]byte(a.Key)) // nolint:errcheck
	// encrypt appid
	h.Write([]byte(shortid.String(a.ID))) // nolint:errcheck
	// encrypt app name
	h.Write([]byte(a.Name)) // nolint:errcheck
	// encrypt uid
	h.Write([]byte(shortid.String(a.Uid))) // nolint:errcheck
	// encrypt timestamp
	h.Write([]byte(ts)) // nolint:errcheck
	return hex.EncodeToString(h.Sum(nil))
}

func (a *App) SetCode() {
	a.Code = RandStr(12)
}

// SetKey must be called after SetCode
func (a *App) SetKey(master string) {
	a.Key = calculateAppKey(*a, master)
}

func (a App) MatchKey(master string) bool {
	return a.Key == calculateAppKey(a, master)
}

// FromBody reads data from request body and unmarshals it to the instance,
// once fields are provided by the caller, an error will be returned if data
// has no such fields or the field values are not set :)
func (a *App) FromBody(r *http.Request, fields ...string) error {
	// TODO: use a faster way to decode json formated request body
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(a); err != nil {
		return err
	}
	if len(fields) != 0 {
		return IsZeroOrInvalid(a, fields...)
	}
	return nil
}

// FromQuery just read fields uid, key and code
func (a *App) FromQuery(r *http.Request, fields ...string) error {
	var err error
	query := r.URL.Query()
	a.ID, err = strconv.Atoi(query.Get("id"))
	if err != nil {
		return auth.ErrQuery
	}
	a.Uid, err = strconv.Atoi(query.Get("uid"))
	if err != nil {
		return auth.ErrQuery
	}
	a.Key = query.Get("key")
	a.Code = query.Get("code")
	if len(fields) != 0 {
		return IsZeroOrInvalid(a, fields...)
	}
	return nil
}

func calculateAppKey(a App, master string) string {
	h := sha256.New()
	h.Write([]byte(master))                // nolint:errcheck
	h.Write([]byte(a.Code))                // nolint:errcheck
	h.Write([]byte(shortid.String(a.ID)))  // nolint:errcheck
	h.Write([]byte(shortid.String(a.Uid))) // nolint:errcheck
	return hex.EncodeToString(h.Sum(nil))
}
