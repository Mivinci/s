package generater

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/handler"
	"github.com/mivinci/s/model"
)

var _ auth.Generater = Generater{}

type Generater struct {
	Key  string
	User *handler.User

	TokenTTL time.Duration
	CodeLen  int
}

func (g Generater) Token(mail string) (token auth.Token, err error) {
	var user model.User
	if err = g.User.DB.One("Mail", mail, &user); err != nil {
		log.Debug("generater.Token: ", err)
		if !errors.Is(err, storm.ErrNotFound) {
			return
		}
		// help unknown user to register
		user = model.User{Mail: mail, Ctime: time.Now(), Name: namegen(12)}
		if err = g.User.DB.Save(&user); err != nil {
			log.Debug("generater.Token: create user error due to ", err)
			return
		}
		log.Debug("register new user: ", user)
	}

	token = auth.Token{
		ID:  user.Uid(), // name is uid for a new user
		Iat: time.Now().Unix(),
		Ein: int64(g.TokenTTL / time.Second),
		Meta: auth.Metadata{
			"uid":   user.ID,
			"name":  user.Name,
			"mail":  user.Mail,
			"perm":  user.Perm,
			"ctime": user.Ctime.Unix(),
		},
		Key: []byte(g.Key),
	}
	log.Debug("generater.Token: generate token: ", token.Final())
	return
}

func (g Generater) Code() (string, error) {
	return codegen(g.CodeLen), nil
}

func codegen(n int) string {
	code := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		code[i] = byte(rand.Int31n(10) + '0')
	}
	return string(code)
}

func namegen(n int) string {
	b := make([]byte, n/2)
	rand.Read(b) // nolint:errcheck
	return fmt.Sprintf("%x", b)
}
