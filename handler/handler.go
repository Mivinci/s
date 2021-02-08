package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mivinci/log"
	"github.com/mivinci/s/render"
)

var (
	ErrAuthHeader = errors.New("no auth token provided by frontier hander")
	ErrHTTPMethod = errors.New("bad http method")
)

var defaultRender = render.Render{Template: retrieveTemplate()}

func JSON(w http.ResponseWriter, data interface{}, err error) { defaultRender.JSON(w, data, err) }

func Error(w http.ResponseWriter, err error) { defaultRender.Error(w, err) }

func OK(w http.ResponseWriter) { defaultRender.OK(w) }

func HTML(w http.ResponseWriter, name string, data interface{}) { defaultRender.HTML(w, name, data) }

func retrieveTemplate() *template.Template {
	var err error
	t := template.New("dnb").Funcs(funcMap)

	if t, err = t.ParseFiles(templateFiles(Env("TEMPLATE_PATH", "html"))...); err != nil {
		log.Fatal("parse templates failed:", err)
	}
	return t
}

func templateFiles(root string) []string {
	ps := make([]string, 0)
	dfs(root, func(path string, fi os.FileInfo) error { // nolint:errcheck
		if filepath.Ext(fi.Name()) == ".html" {
			ps = append(ps, path)
			log.Debug("find template:", path)
		}
		return nil
	})
	return ps
}

var funcMap = template.FuncMap{
	"upper":     strings.ToUpper,
	"lower":     strings.ToLower,
	"trimLeft":  strings.TrimLeft,
	"trimRight": strings.TrimRight,
	"contains":  strings.Contains,
	"duration":  duration,
}

func duration(t time.Time) string {
	d := time.Since(t)
	if d < time.Minute {
		return fmt.Sprintf("%d秒前", d/time.Second)
	}
	if d < time.Hour {
		return fmt.Sprintf("%d分钟前", d/time.Minute)
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d小时前", d/time.Hour)
	}
	return t.Format("2006-01-02")
}

func dfs(root string, fn func(path string, fi os.FileInfo) error) error {
	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		path := filepath.Join(root, fi.Name())
		if err := fn(path, fi); err != nil {
			return err
		}
		if fi.IsDir() {
			if err := dfs(path, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func Env(key, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

func randStr(n int) string {
	b := make([]byte, n/2)
	rand.Read(b) // nolint:errcheck
	return fmt.Sprintf("%x", b)
}

func randStr12() string { return randStr(12) }
