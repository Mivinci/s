package filter

import (
	"net/http"

	"github.com/mivinci/log"
	"github.com/mivinci/s/auth"
	"github.com/mivinci/s/model"
)

func UnmarshalApp(r *http.Request, master string, app *model.App, fields ...string) error {
	uid, ok := FromConext(r)
	if !ok {
		return auth.ErrPermDenied
	}

	// decode request and unmarshal app
	if r.Method == "GET" {
		if err := app.FromQuery(r, fields...); err != nil {
			return err
		}
	} else {
		if err := app.FromBody(r, fields...); err != nil {
			return err
		}
	}

	log.Debug("filter.UnmarshalApp: ", app)

	//
	if app.Uid != uid {
		return auth.ErrPermDenied
	}

	if master != "" {
		// avoid access from user to app of other users.
		if !app.MatchKey(master) {
			return auth.ErrAccessForbidden
		}
	}

	return nil
}
