package api

import (
	"net/http"

	"github.com/euforia/thrap/thrapb"

	"github.com/gorilla/mux"
)

func (api *httpHandler) handleListDeployments(w http.ResponseWriter, r *http.Request) {
	projID := mux.Vars(r)["pid"]
	proj, err := api.projs.Get(projID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	d := proj.Deployments()
	deploys, err := d.List()

	writeJSONResponse(w, deploys, err)
}

func (api *httpHandler) handleDeployment(w http.ResponseWriter, r *http.Request) {
	var (
		vars   = mux.Vars(r)
		projID = vars["pid"]
		envID  = vars["eid"]
		instID = vars["iid"]
		resp   interface{}
	)

	proj, err := api.projs.Get(projID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	dpl := proj.Deployments()

	switch r.Method {
	case "GET":
		resp, err = dpl.Get(envID, instID)

	case "POST":
		_, err = dpl.Create(&thrapb.Deployment{
			Name:    instID,
			Profile: &thrapb.Profile{ID: envID},
		})

	// case "PUT":
	// 	err=dpl.Deploy(depl)

	case http.MethodOptions:
	default:
		w.WriteHeader(405)
		return
	}

	// w.WriteHeader(404)
	// w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")

	writeJSONResponse(w, resp, err)
}
