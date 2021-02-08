package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mivinci/log"
)

var (
	ErrTokenDecode      = errors.New("token cannot be decoded")
	ErrTokenNotmatch    = errors.New("token does not match")
	ErrTokenExpired     = errors.New("token is dead")
	ErrTokenIncompleted = errors.New("incompleted token field")
)

type Token struct {
	ID   string   `json:"id"`
	Iat  int64    `json:"iat"`
	Ein  int64    `json:"ein"`
	Meta Metadata `json:"meta"`
	Key  []byte   `json:"key,omitempty"`
}

func (t Token) Final() string {
	return fmt.Sprintf("%s.%s", t.Encode(), t.Checksum())
}

func (t Token) Encode() string {
	t.Key = nil
	buf, _ := json.Marshal(&t)
	return base64.URLEncoding.EncodeToString(buf)
}

func (t *Token) Decode(s string) error {
	buf, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, t)
}

func (t Token) Checksum() string {
	h := sha256.New()
	b, _ := json.Marshal(&t)
	h.Write(b)     // nolint:errcheck
	h.Write(t.Key) // nolint:errcheck
	return hex.EncodeToString(h.Sum(nil))
}

func (t Token) Match(sum string) bool {
	return sum == t.Checksum()
}

func (t *Token) OK(s string) error {
	ss := strings.Split(s, ".")
	if t.Decode(ss[0]) != nil {
		return ErrTokenDecode
	}

	log.Debug("decoded: ", t)

	if !t.Match(ss[1]) {
		return ErrTokenNotmatch
	}
	if t.Iat+t.Ein < time.Now().Unix() {
		return ErrTokenExpired
	}
	return nil
}

func (t *Token) OKUnless(s string, fn func(*Token) error) error {
	if err := t.OK(s); err != nil {
		return err
	}
	return fn(t)
}
