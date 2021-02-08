package render

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/mivinci/log"
)

type protocol struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

type Render struct {
	Template *template.Template
}

func (Render) JSON(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	var errmsg string
	if err != nil {
		errmsg = err.Error()
	}
	b, err := json.Marshal(&protocol{
		Code:  0,
		Error: errmsg,
		Data:  data,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(b) // nolint:errcheck
}

func (r Render) Error(w http.ResponseWriter, err error) { r.JSON(w, nil, err) }

func (r Render) OK(w http.ResponseWriter) { r.JSON(w, "ok", nil) }

func (r Render) HTML(w http.ResponseWriter, name string, data interface{}) {
	if err := r.Template.ExecuteTemplate(w, name, data); err != nil {
		log.Error("render.HTML:", err)
	}
}
