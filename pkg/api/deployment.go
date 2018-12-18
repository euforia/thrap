package api

import (
	"encoding/json"
	"net/http"

	"github.com/euforia/thrap/pkg/thrap"

	"github.com/gorilla/mux"
)

func (api *httpHandler) handleListDeployments(w http.ResponseWriter, r *http.Request) {
	projID := mux.Vars(r)["pid"]
	proj, err := api.projects.Get(projID)
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

		d    *thrap.Deployment
		resp interface{}
	)

	proj, err := api.projects.Get(projID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	dpl := proj.Deployments()
	ctx := r.Context()

	switch r.Method {

	case "GET":
		d, err = dpl.Get(ctx, envID, instID)

	case "PUT":
		d, err = dpl.Create(ctx, envID, instID)

	// case "POST":
	// 	defer r.Body.Close()

	// 	d, err = dpl.Get(ctx, envID, instID)
	// 	if err == nil {
	// 		err = api.handleDeploy(d, r)
	// 	}

	// case "DELETE":
	// 	// stop and purge
	// 	d, err = dpl.Get(ctx, envID, instID)
	// 	if err == nil {
	// 		err = d.Destroy(ctx)
	// 	}

	case http.MethodOptions:
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT")
		w.WriteHeader(200)
		return

	default:
		w.WriteHeader(405)
		return
	}

	if err == nil {
		resp = d.Deployable()
	}

	writeJSONResponse(w, resp, err)
}

func (api *httpHandler) handleDeploy(w http.ResponseWriter, r *http.Request) {
	var (
		vars   = mux.Vars(r)
		projID = vars["pid"]
		profID = vars["eid"]
		instID = vars["iid"]

		d *thrap.Deployment

		err  error
		resp interface{}

		ctx = r.Context()
	)

	proj, err := api.projects.Get(projID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	dpl := proj.Deployments()
	d, err = dpl.Get(ctx, profID, instID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	switch r.Method {
	case "POST":
		dec := json.NewDecoder(r.Body)
		var req thrap.DeployRequest
		err = dec.Decode(&req)
		if err == nil {
			resp, err = d.Deploy(ctx, &req)
		}

	case "DELETE":
		if _, purge := r.URL.Query()["purge"]; purge {
			err = d.Destroy(ctx)
		} else {
			err = d.Stop(ctx)
		}

		if err == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

	case http.MethodOptions:
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "POST,DELETE")
		w.WriteHeader(200)
		return

	default:
		w.WriteHeader(405)
		return
	}

	writeJSONResponse(w, resp, err)
}
