package api

import (
	"net/http"

	"github.com/euforia/thrap/thrapb"
)

func (api *httpHandler) handleListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles := api.t.Profiles()
	list := profiles.List()

	writeJSONResponse(w, list, nil)
}

func (api *httpHandler) handleProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		resp *thrapb.Project
		err  error
	)

	switch r.Method {
	case "GET":
		resp, err = api.getProject(w, r)

	case "OPTIONS":
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.WriteHeader(200)
		return

	default:
		w.WriteHeader(405)
		return
	}

	writeJSONResponse(w, resp, err)
}
