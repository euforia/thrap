package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/pkg/thrap"
	"github.com/gorilla/mux"
)

func (api *httpHandler) setDeploymentSpec(r *http.Request, dpl *thrap.Deployments, version string) ([]byte, error) {
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var (
		desc        *pb.DeploymentDescriptor
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
		desc.ID = version
		if err = dpl.SetDescriptor(desc); err == nil {
			return desc.Spec, nil
		}
	}

	return nil, err
}

func (api *httpHandler) handleListDeploymentSpecs(w http.ResponseWriter, r *http.Request) {
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
		var descs []string
		descs, err = dpl.Descriptors()
		if err == nil {
			resp, err = json.Marshal(descs)
		}

	case http.MethodOptions:
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET")
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

func (api *httpHandler) handleDeploymentSpec(w http.ResponseWriter, r *http.Request) {
	projID := mux.Vars(r)["pid"]
	proj, err := api.projects.Get(projID)
	if err != nil {
		err = errors.Wrapf(err, "project %s", projID)
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	version := mux.Vars(r)["version"]
	if len(version) == 0 {
		w.WriteHeader(404)
		return
	}

	var (
		dpl  = proj.Deployments()
		resp []byte
	)

	switch r.Method {

	case "GET":
		desc, err := dpl.Descriptor(version)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte(err.Error()))
			return
		}
		resp = desc.Spec

	case "PUT":
		resp, err = api.setDeploymentSpec(r, dpl, version)

	case "DELETE":
		err = dpl.DeleteDescriptor(version)

	case http.MethodOptions:
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,DELETE")
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
