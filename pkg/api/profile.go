package api

import (
	"net/http"

	"github.com/euforia/thrap/pkg/storage"
	"github.com/gorilla/mux"
)

func (api *httpHandler) handleListProfiles(w http.ResponseWriter, r *http.Request) {
	profiles := api.t.Profiles()
	list := profiles.List()

	writeJSONResponse(w, list, nil)
}

func (api *httpHandler) handleProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "GET":
		profs := api.t.Profiles()
		prof := mux.Vars(r)["id"]
		resp, err := profs.Get(prof)
		if err == storage.ErrProfileNotFound {
			w.WriteHeader(404)
			return
		}
		writeJSONResponse(w, resp, err)

	case "OPTIONS":
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.WriteHeader(200)
		return

	default:
		w.WriteHeader(405)
		return
	}

}
