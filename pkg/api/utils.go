package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/euforia/thrap/thrapb"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
)

func setAccessControlHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func writeJSONResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if resp == nil {
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	setAccessControlHeaders(w)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// JSON Marshal nomad job, store in descriptor. It returns aDeploymentDescriptor
func makeNomadJSONDeployDesc(job *nomad.Job) (*thrapb.DeploymentDescriptor, error) {
	nb, err := json.Marshal(job)
	if err != nil {
		return nil, err
	}

	return &thrapb.DeploymentDescriptor{
		Spec: nb,
		Mime: DescContentTypeJSON,
	}, nil
}

func parseNomadHCLDescriptor(b []byte) (*thrapb.DeploymentDescriptor, error) {
	// var out map[string]interface{}
	// err := hcl.Unmarshal(b, &out)
	job, err := jobspec.Parse(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	return makeNomadJSONDeployDesc(job)
}
