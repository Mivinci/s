package auth

import (
	"github.com/mivinci/log"
)

type Sender interface {
	Send(to, code string) error
}

// sender implements Sender for debuging
type sender struct{}

func (sender) Send(to, code string) error {
	log.Debugf("send code '%s' to %s\n", code, to)
	return nil
}
