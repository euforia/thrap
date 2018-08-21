package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/thrap"
	"github.com/euforia/thrap/thrapb"
	"github.com/gorilla/mux"
)

func (api *httpHandler) handleDeploymentSpec(w http.ResponseWriter, r *http.Request) {
	projID := mux.Vars(r)["pid"]
	proj, err := api.projects.Get(projID)
	if err != nil {
		err = errors.Wrapf(err, "project %s", projID)

		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	var (
		dpl  = proj.Deployments()
		resp []byte
	)

	switch r.Method {
	case "GET":
		desc := dpl.Descriptor()
		if len(desc.Spec) == 0 {
			w.WriteHeader(404)
			return
		}
		resp = desc.Spec

	case "POST":
		resp, err = api.setDeploymentSpec(r, dpl)

	default:
		w.WriteHeader(405)
		return
	}

	if err != nil {
		w.WriteHeader(400)
		setAccessControlHeaders(w)
		w.Write([]byte(err.Error()))
	} else {
		setAccessControlHeaders(w)
		w.Write(resp)
	}
}

func (api *httpHandler) setDeploymentSpec(r *http.Request, dpl *thrap.Deployments) ([]byte, error) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var desc thrapb.DeploymentDescriptor

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/hcl":
		var out map[string]interface{}
		err = hcl.Unmarshal(b, &out)
		if err == nil {
			desc = thrapb.DeploymentDescriptor{
				Spec: b,
			}
		}
	case "application/json":
		err = json.Unmarshal(b, &desc)

	default:
		err = fmt.Errorf("unknown content-type: '%s'", contentType)

	}

	if err == nil {
		err = dpl.SetDescriptor(&desc)
	}

	return desc.Spec, err
}
