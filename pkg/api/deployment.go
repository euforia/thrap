package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/euforia/thrap/pkg/thrap"

	"github.com/gorilla/mux"
)

// swagger:operation GET /project/{projectId}/deployments listDeployments
//
// Returns a list of project deployments
//
// Returns a list of project deployments
//
// ---
// responses:
//   200:
//     description: "OK"
//     schema:
//       type: array
//   404:
//     description: "project not found"
//   500:
//     description: "Internal Server Error"
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
	// swagger:operation GET /project/{projectId}/deployment/{environmentId}/{instanceId} getDeployment
	//
	// Returns information about a deployment
	//
	// Returns information about a deployment
	//
	// ---
	// responses:
	//   200:
	//     description: "OK"
	//   404:
	//     description: "deployment not found"
	//   500:
	//     description: "Internal Server Error"
	case "GET":
		d, err = dpl.Get(ctx, envID, instID)

		// swagger:operation POST /project/{projectId}/deployment/{environmentId}/{instanceId} createDeployment
		//
		// Create a new deployment for a project
		//
		// Create a new deployment for a project
		//
		// ---
		// responses:
		//   200:
		//     description: "OK"
		//   400:
		//     description: "failed to create deployment"
		//   500:
		//     description: "Internal Server Error"
	case "POST":
		d, err = dpl.Create(ctx, envID, instID)

		// swagger:operation PUT /project/{projectId}/deployment/{environmentId}/{instanceId} deploy
		//
		// Deploy
		//
		// Deploy
		//
		// ---
		// responses:
		//   200:
		//     description: "OK"
		//   400:
		//     description: "failed to deploy"
		//   500:
		//     description: "Internal Server Error"
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

	_, err = d.Deploy(&thrap.DeployRequest{Variables: vars})
	return err
}
