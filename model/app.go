package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/mivinci/shortid"
)

type App struct {
	ID    int       `json:"id" storm:"increment"`
	Uid   int       `json:"uid" storm:"unique"`
	Name  string    `json:"name" storm:"unique"`
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
