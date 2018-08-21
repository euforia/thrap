package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/euforia/thrap/pkg/thrap"
	"github.com/euforia/thrap/thrapb"
	"github.com/gorilla/mux"
)

func (api *httpHandler) handleDeploymentSpec(w http.ResponseWriter, r *http.Request) {
	projID := mux.Vars(r)["pid"]
	proj, err := api.projects.Get(projID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	var (
		dpl  = proj.Deployments()
		resp interface{}
	)

	switch r.Method {
	case "GET":
		resp = dpl.Descriptor()

	case "POST":
		err = api.setDeploymentSpec(r, dpl)

	default:
		w.WriteHeader(405)
		return
	}

	writeJSONResponse(w, resp, err)
}

func (api *httpHandler) setDeploymentSpec(r *http.Request, dpl *thrap.Deployments) error {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var spec thrapb.DeploymentDescriptor
	err = json.Unmarshal(b, &spec)
	if err != nil {
		return err
	}

	return dpl.SetDescriptor(&spec)
}
