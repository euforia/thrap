package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/euforia/thrap/pkg/thrap"
	"github.com/euforia/thrap/thrapb"

	"github.com/gorilla/mux"
)

// swagger:operation GET /projects listProjects
//
// Returns a list of projects
//
// Returns a list of projects
//
// ---
// responses:
//   200:
//     description: "OK"
//     schema:
//       type: array
//   400:
//     description: "Bad Request"
//   500:
//     description: "Internal Server Error"
func (api *httpHandler) handleListProjects(w http.ResponseWriter, r *http.Request) {
	list := make([]*thrapb.Project, 0)

	err := api.projects.Iter("", func(proj *thrapb.Project) error {
		list = append(list, proj)
		return nil
	})

	writeJSONResponse(w, list, err)
}

func (api *httpHandler) handleProject(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		resp *thrapb.Project
		err  error
	)

	switch r.Method {
	// swagger:operation GET /project/{projectId} getProject
	//
	// Get project details
	//
	// Get project details
	//
	// ---
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
		resp, err = api.getProject(w, r)

	// swagger:operation POST /project/{projectId} createProject
	//
	// Create a new project
	//
	// Create a new project
	//
	// ---
	// responses:
	//   200:
	//     description: "OK"
	//   400:
	//     description: "Bad Request"
	//   500:
	//     description: "Internal Server Error"
	case "POST":
		resp, err = api.createProject(w, r)

		// swagger:operation PUT /project/{projectId} updateProject
		//
		// Update a project
		//
		// Update a project
		//
		// ---
		// responses:
		//   200:
		//     description: "OK"
		//   400:
		//     description: "Bad Request"
		//   500:
		//     description: "Internal Server Error"
	case "PUT":
		resp, err = api.updateProject(w, r)

	// case "DELETE":
	// 	err = api.delete(w, r)

	case "OPTIONS":
		setAccessControlHeaders(w)
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT")
		w.WriteHeader(200)
		return

	default:
		w.WriteHeader(405)
		return
	}

	writeJSONResponse(w, resp, err)
}

func (api *httpHandler) getProject(w http.ResponseWriter, r *http.Request) (*thrapb.Project, error) {
	projID := mux.Vars(r)["id"]
	proj, err := api.projects.Get(projID)
	if err == nil {
		return proj.Project, nil
	}
	return nil, err
}

func (api *httpHandler) createProject(w http.ResponseWriter, r *http.Request) (*thrapb.Project, error) {
	projID := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var req thrap.ProjectCreateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}

	if req.Project == nil {
		req.Project = &thrapb.Project{ID: projID}
	} else {
		req.Project.ID = projID
	}

	ctx := r.Context()
	nProj, err := api.projects.Create(ctx, &req)
	if err == nil {
		return nProj.Project, nil
	}
	return nil, err
}

func (api *httpHandler) updateProject(w http.ResponseWriter, r *http.Request) (*thrapb.Project, error) {
	projID := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var proj thrapb.Project
	err = json.Unmarshal(body, &proj)
	if err != nil {
		return nil, err
	}
	proj.ID = projID

	sproj, err := api.projects.Get(projID)
	if err != nil {
		return nil, err
	}
	sproj.Project = &proj
	err = sproj.Sync()

	return &proj, err
}

func (api *httpHandler) delete(w http.ResponseWriter, r *http.Request) error {
	projID := mux.Vars(r)["id"]
	return api.projects.Delete(projID)
}
