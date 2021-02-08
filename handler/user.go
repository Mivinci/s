package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/mivinci/log"
	"github.com/mivinci/s/model"
)

var ErrParam = errors.New("invalid post parameter")

type User struct {
	DB *storm.DB
}

func (u *User) Init(name string) error {
	return u.DB.Save(&model.User{Name: name, Ctime: time.Now(), Mail: "unknown"})
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	name, email := r.PostFormValue("name"), r.PostFormValue("email")
	if name == "" || email == "" {
		Error(w, ErrParam)
		return
	}

	user := model.User{Name: name, Mail: email, Ctime: time.Now()}
	if err := u.DB.Save(&user); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	JSON(w, &user, nil)
}

func (u *User) One(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		Error(w, ErrParam)
		return
	}
	var user model.User
	if err := u.DB.One("ID", id, &user); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	JSON(w, &user, nil)
}

func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		Error(w, ErrParam)
		return
	}
	name, email := r.PostFormValue("name"), r.PostFormValue("email")
	user := model.User{ID: id, Name: name, Mail: email}
	if err := u.DB.Update(&user); err != nil {
		log.Error(err)
		Error(w, err)
		return
	}
	OK(w)
}

// Delete hard delete
func (u *User) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		Error(w, ErrParam)
		return
	}
	if err := u.DB.DeleteStruct(&model.User{ID: id}); err != nil {
		log.Error(err)
		Error(w, ErrParam)
	}
	OK(w)
}

func (u *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u.One(w, r)
	case "POST":
		u.Create(w, r)
	case "PUT":
		u.Update(w, r)
	case "DELETE":
		u.Delete(w, r)
	default:
		Error(w, ErrHTTPMethod)
	}
}
