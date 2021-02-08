package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/mivinci/log"
	"github.com/mivinci/s/model"
)

type App struct {
	DB *storm.DB
}

func (a *App) Init(name string) error {
	return a.DB.Save(&model.App{Name: name, Uid: 1, Ctime: time.Now()})
}

func (a *App) One(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		Error(w, ErrParam)
		return
	}
	var ap model.App
	if err := a.DB.One("ID", id, &ap); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	JSON(w, &ap, nil)
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	uid, err := strconv.Atoi(r.PostFormValue("uid"))
	name, site := r.PostFormValue("name"), r.PostFormValue("site")
	log.Debugf("uid(%d) name(%s) site(%s)\n", uid, name, site)
	if err != nil || name == "" || site == "" {
		Error(w, ErrParam)
		return
	}
	key := randStr12()
	ap := model.App{Uid: uid, Name: name, Site: site, Ctime: time.Now(), Key: key}
	if err := a.DB.Save(&ap); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	JSON(w, &ap, nil)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		Error(w, ErrParam)
		return
	}
	name, site := r.PostFormValue("name"), r.PostFormValue("site")
	ap := model.App{ID: id, Name: name, Site: site}
	if err := a.DB.Update(&ap); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	OK(w)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		Error(w, ErrParam)
		return
	}
	if err := a.DB.DeleteStruct(&model.App{ID: id}); err != nil {
		log.Error(err)
		Error(w, ErrParam)
	}
	OK(w)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.One(w, r)
	case "POST":
		a.Create(w, r)
	case "PUT":
		a.Update(w, r)
	case "DELETE":
		a.Delete(w, r)
	default:
		Error(w, ErrHTTPMethod)
	}
}
