package handler

import (
	"encoding/json"
	"errors"
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
	var app model.App
	err := FromBody(r, &app)
	if err != nil {
		log.Debug("decode from request body: ", err)
		Error(w, ErrParam)
		return
	}
	if app.Site == "" {
		Error(w, errors.New("app missing field `site`"))
		return
	}
	app.Ctime = time.Now()
	app.Key = randStr12()
	log.Debug("create app: ", app)
	if err := a.DB.Save(&app); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	JSON(w, &app, nil)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	var app model.App
	err := FromBody(r, &app)
	if err != nil {
		log.Debug("decode from request body: ", err)
		Error(w, err)
		return
	}
	if err := a.DB.Update(&app); err != nil {
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

	// TODO: check if the user owns the app
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

func FromBody(r *http.Request, app *model.App) error {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := dec.Decode(app); err != nil {
		return err
	}
	// app.ID = 0
	app.Key = ""
	return nil
}
