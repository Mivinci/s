package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/asdine/storm/v3"
	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/model"
)

type Template struct {
	SubTitle string
	Notice   Notice
	User     model.User
}

type Notice struct {
	Type, Text string
}

type AuthTemplte struct {
	Template
	App     model.App
	CodeLen int
}

type ConsoleTemplate struct {
	Template
	User model.User
	Apps []model.App
}

type ErrorTemplate struct {
	Template
	Error string
}

type View struct {
	App     *App
	User    *User
	CodeLen int
}

func (v *View) Auth(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	appid, err := strconv.Atoi(query.Get("appid"))
	sign, ts := query.Get("sign"), query.Get("ts")

	t := Template{SubTitle: "登录"}

	if err != nil || sign == "" || ts == "" {
		HTML(w, "error", &ErrorTemplate{Template: t, Error: "no appid provided"})
		return
	}

	// TODO: check if ts is in a period of now

	var app model.App
	if err := v.App.DB.One("ID", appid, &app); err != nil {
		log.Error(appid, err)
		HTML(w, "error", &ErrorTemplate{Template: t, Error: "unknown application"})
		return
	}

	if sum := app.Checksum(ts); sum != sign {
		log.Debug("app:", app, " sum:", sum)
		HTML(w, "error", &ErrorTemplate{Template: t, Error: "incorrect signature"})
		return
	}

	HTML(w, "auth", &AuthTemplte{Template: t, App: app, CodeLen: v.CodeLen})
}

func (v *View) Console(w http.ResponseWriter, r *http.Request) {
	t := Template{SubTitle: "用户"}

	md, ok := auth.FromContext(r)
	if !ok {
		HTML(w, "error", &ErrorTemplate{Template: t, Error: auth.ErrMetadata.Error()})
		return
	}

	uid, ok := md["uid"].(int)
	if !ok {
		HTML(w, "error", &ErrorTemplate{Template: t, Error: "unknown uid"})
		return
	}

	var user model.User
	if err := v.User.DB.One("ID", uid, &user); err != nil {
		log.Error(err)
		HTML(w, "error", &ErrorTemplate{Template: t, Error: "unknown model.User"})
		return
	}

	var apps []model.App
	if err := v.App.DB.Find("Uid", uid, &apps); err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			log.Error(err)
			HTML(w, "error", &ErrorTemplate{Template: t, Error: err.Error()})
			return
		}
	}

	t.SubTitle = user.Name
	t.User = user
	HTML(w, "console", &ConsoleTemplate{Template: t, User: user, Apps: apps})
}

func (v *View) Edit() {}
