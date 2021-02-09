package main

import (
	"net/http"
	"os"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/gob"
	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/filter"
	"github.com/mivinci/s/generater"
	"github.com/mivinci/s/handler"
	"github.com/mivinci/s/render"
	"github.com/mivinci/s/sender"
)

var (
	tokenTTL = time.Hour * 24 * 7
	codeTTL  = time.Minute * 10
	codeLen  = 4
)

func init() {
	logMode := Env("MODE", "development")
	if logMode == "production" {
		log.SetMode(log.ModeProduct)
	}
}

func main() {
	path := Env("DB", "data/db")
	key := Env("KEY", "wfs1101")
	owner := Env("OWNER", "绿椰子")
	self := "绿椰子"

	db, err := storm.Open(path, storm.Codec(gob.Codec))
	if err != nil {
		log.Fatal("open bolt:", err)
	}

	app := handler.App{DB: db, Key: key}
	user := handler.User{DB: db}
	view := handler.View{App: &app, User: &user, CodeLen: codeLen}

	e1, e2 := user.Init(owner), app.Init(self)
	if e1 != nil || e2 != nil {
		log.Warn("user ", e1, ", app ", e2)
	}

	aut := auth.New(key,
		auth.WithCodeTTL(codeTTL),
		auth.WithTokenTTL(tokenTTL),
		auth.WithSender(sender.Sender{}),
		auth.WithRender(render.Render{}),
		auth.WithFilter(filter.Filter{}),
		auth.WithGenerater(generater.Generater{Key: key, User: &user, TokenTTL: tokenTTL, CodeLen: codeLen}),
	)

	tp := handler.Transport{Auth: aut}

	// statics
	http.HandleFunc("/favicon.ico", faviconHander)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	// auth api
	http.HandleFunc("/transport", aut.Transport(&tp, "GET"))

	// curd api
	http.Handle("/app", aut.Private(&app))
	http.Handle("/user", aut.Private(&user))

	// views
	http.HandleFunc("/auth", view.Auth)
	http.HandleFunc("/console", aut.PrivateFunc(view.Console))

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("listen:", err)
	}
}

func Env(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

func faviconHander(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}
