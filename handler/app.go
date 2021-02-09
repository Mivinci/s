package handler

import (
	"net/http"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/mivinci/log"
	"github.com/mivinci/s/filter"
	"github.com/mivinci/s/model"
)

type App struct {
	DB  *storm.DB
	Key string
}

func (a *App) Init(name string) error {
	if err := a.DB.Init(model.App{}); err != nil {
		return err
	}
	return a.DB.Save(&model.App{Name: name, Uid: 1, Ctime: time.Now(), Type: model.TypeWeb, State: model.StateOK})
}

func (a *App) One(w http.ResponseWriter, r *http.Request) {
	app := new(model.App)
	err := filter.UnmarshalApp(r, a.Key, app, "ID", "Uid", "Key", "Code")
	if err != nil {
		log.Debug("App.One: ", err)
		Error(w, err)
		return
	}
	if err := a.DB.One("ID", app.ID, app); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	JSON(w, app, nil)
}

func (a *App) Create(w http.ResponseWriter, r *http.Request) {
	app := new(model.App)
	err := filter.UnmarshalApp(r, "", app, "Site")

	if err != nil {
		log.Debug("App.Create: ", err)
		Error(w, err)
		return
	}

	app.Ctime = time.Now()
	app.SetCode()
	if err := a.DB.Save(app); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	app.SetKey(a.Key)
	log.Debug("App.Create: ", app)
	JSON(w, app, nil)
}

func (a *App) Update(w http.ResponseWriter, r *http.Request) {
	app := new(model.App)
	err := filter.UnmarshalApp(r, a.Key, app, "ID", "Uid", "Key", "Code")

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
	app := new(model.App)
	err := filter.UnmarshalApp(r, a.Key, app, "ID", "Uid", "Key", "Code")

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
