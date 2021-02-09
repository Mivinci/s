package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mivinci/log"
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
	Key   string    `json:"key"`
	Ctime time.Time `json:"ctime" storm:"index"`
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

func (a *App) SetKey(master string) {
	a.Key = CalculateAppKey(a, master)
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
	a.Key = "" // avoid setting key by caller
	if len(fields) != 0 {
		return IsZeroOrInvalid(a, fields...)
	}
	return nil
}

// AccessAndRead reads app info from request and see if the user has access to it by:
// 1. checking if the uid extracted from token is the same as the uid from the app info
// 2. comparing the received key with the secret key of the app
func (a *App) AccessAndRead(r *http.Request, uid int, master string, fields ...string) error {
	log.Debug("model.App.AccessAndRead: uid ", uid)
	err := a.FromBody(r, fields...)
	if err != nil {
		return err
	}

	if a.Uid != uid {
		return auth.ErrPermDenied
	}

	if a.Key != CalculateAppKey(a, master) {
		return auth.ErrAccessForbidden
	}

	return nil
}

func CalculateAppKey(app *App, master string) string {
	h := sha256.New()
	h.Write([]byte(master))                  // nolint:errcheck
	h.Write([]byte(shortid.String(app.ID)))  // nolint:errcheck
	h.Write([]byte(shortid.String(app.Uid))) // nolint:errcheck
	return hex.EncodeToString(h.Sum(nil))
}
