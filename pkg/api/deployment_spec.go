package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

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

	// swagger:operation GET /project/{project}/deployment/spec getDeploySpec
	//
	// Get deployment spec for a project
	//
	// Get deployment spec for a project
	//
	// responses:
	//   200:
	//     description: "OK"
	//   400:
	//     description: "Bad Request"
	//   404:
	//     description: "specification not found"
	//   500:
	//     description: "Internal Server Error"
	case "GET":
		desc := dpl.Descriptor()
		if desc == nil || len(desc.Spec) == 0 {
			w.WriteHeader(404)
			return
		}
		resp = desc.Spec

		// swagger:operation POST /project/{project}/deployment/spec updateDeploySpec
		//
		// Update deployment spec for a project
		//
		// Update deployment spec for a project
		//
		// responses:
		//   200:
		//     description: "OK"
		//   400:
		//     description: "Bad Request"
		//   500:
		//     description: "Internal Server Error"
	case "POST":
		resp, err = api.setDeploymentSpec(r, dpl)

	case http.MethodOptions:
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		w.WriteHeader(200)
		return

	default:
		w.WriteHeader(405)
		return
	}

	if err != nil {
		setAccessControlHeaders(w)
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		setAccessControlHeaders(w)
		w.WriteHeader(200)
		w.Write(resp)
	}
}

func (api *httpHandler) setDeploymentSpec(r *http.Request, dpl *thrap.Deployments) ([]byte, error) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var (
		desc        *thrapb.DeploymentDescriptor
		contentType = r.Header.Get("Content-Type")
	)

	switch contentType {

	case DescContentTypeMoldHCL:
		desc, err = parseMoldHCLDeployDesc(b)

	case DescContentTypeNomadHCL:
		desc, err = parseNomadHCLDescriptor(b)

	case DescContentTypeNomadJSON:
		desc, err = parseNomadJSONDescriptor(b)

	default:
		err = fmt.Errorf("unsupported Content-Type: '%s'", contentType)

	}

	if err == nil {
		if err = dpl.SetDescriptor(desc); err == nil {
			return desc.Spec, nil
		}
	}

	return nil, err
}
