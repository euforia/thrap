package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/euforia/thrap/pkg/storage"
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

	d := proj.Deployments(storage.DefaultSpecVersion)
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

	dpl := proj.Deployments(storage.DefaultSpecVersion)
	ctx := r.Context()

	switch r.Method {

	case "GET":
		d, err = dpl.Get(ctx, envID, instID)

	case "POST":
		d, err = dpl.Create(ctx, envID, instID)

	case "PUT":
		defer r.Body.Close()

		d, err = dpl.Get(ctx, envID, instID)
		if err == nil {
			err = api.handleDeploy(d, r)
		}

	case http.MethodOptions:
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT")
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

func (api *httpHandler) handleDeploy(d *thrap.Deployment, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var vars map[string]string
	if err = json.Unmarshal(b, &vars); err != nil {
		return err
	}

	_, err = d.Deploy(r.Context(), &thrap.DeployRequest{Variables: vars})
	return err
}
