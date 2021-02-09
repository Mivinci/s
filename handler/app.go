package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/filter"
	"github.com/mivinci/s/model"
)

type App struct {
	DB  *storm.DB
	Key string
}

func (a *App) Init(name string) error {
	return a.DB.Save(&model.App{Name: name, Uid: 1, Ctime: time.Now(), Type: model.TypeWeb, State: model.StateOK})
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
	uid, ok := filter.FromConext(r)
	if !ok {
		Error(w, auth.ErrPermDenied)
		return
	}

	app := new(model.App)
	err := app.AccessAndRead(r, uid, "Site")
	if err != nil {
		log.Debug("App.Update: ", err)
		Error(w, err)
		return
	}

	app.Ctime = time.Now()
	log.Debug("App.Create: ", app)
	// if err := a.DB.Save(&app); err != nil {
	// 	log.Error(err)
	// 	Error(w, err)
	// 	return
	// }
	app.SetKey(a.Key)

	JSON(w, &app, nil)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	uid, ok := filter.FromConext(r)
	if !ok {
		Error(w, auth.ErrPermDenied)
		return
	}

	app := new(model.App)
	err := app.AccessAndRead(r, uid, "ID", "Uid", "Key")
	if err != nil {
		log.Debug("App.Update: ", err)
		Error(w, err)
		return
	}

	if err := a.DB.Update(app); err != nil {
		log.Error("App.Update: ", err)
		Error(w, err)
		return
	}
	OK(w)
}

func (a *App) Delete(w http.ResponseWriter, r *http.Request) {
	uid, ok := filter.FromConext(r)
	if !ok {
		Error(w, auth.ErrPermDenied)
		return
	}

	app := new(model.App)
	err := app.AccessAndRead(r, uid, "ID", "Uid", "Key")
	if err != nil {
		log.Debug("App.Update: ", err)
		Error(w, err)
		return
	}

	// TODO: check if the user owns the app
	if err := a.DB.DeleteStruct(app); err != nil {
		log.Error(err)
		Error(w, ErrParam)
		return
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
