package api

import (
	"net"
	"net/http"

	"github.com/euforia/thrap/pkg/thrap"
	"github.com/gorilla/mux"
)

const (
	DescContentTypeMold  = "application/vnd.thrap.mold.deployment.descriptor.v1+hcl"
	DescContentTypeNomad = "application/vnd.thrap.nomad.deployment.descriptor.v1+hcl"
	DescContentTypeJSON  = "application/json"
)

// Config is the api server config
type Config struct {
	// Projects *project.Projects
	// Profiles storage.ProfileStorage
}

type httpHandler struct {
	t        *thrap.Thrap
	projects *thrap.Projects
}

// Server is the REST api interface
type Server struct {
	// api request router
	router *mux.Router
	// api handler
	handler *httpHandler
}

// NewServer returns a new API server
func NewServer(t *thrap.Thrap) *Server {
	server := &Server{
		router: mux.NewRouter(),
		handler: &httpHandler{
			t:        t,
			projects: thrap.NewProjects(t),
		},
	}

	server.registerHandlers()

	return server
}

func (server *Server) registerHandlers() {
	// server.router.HandleFunc("/v1/identities", server.ident.list)
	// server.router.HandleFunc("/v1/identity/{id}", server.ident.identity)

	server.router.HandleFunc("/v1/profiles", server.handler.handleListProfiles)
	server.router.HandleFunc("/v1/profile/{id}", server.handler.handleProfile)

	server.router.HandleFunc("/v1/projects", server.handler.handleListProjects)
	server.router.HandleFunc("/v1/project/{id}", server.handler.handleProject)

	server.router.HandleFunc("/v1/project/{pid}/deployments", server.handler.handleListDeployments)
	server.router.HandleFunc("/v1/project/{pid}/deployment/spec", server.handler.handleDeploymentSpec)
	server.router.HandleFunc("/v1/project/{pid}/deployment/{eid}/{iid}", server.handler.handleDeployment)
}

// Serve starts serving the registered handlers on the given listener
func (server *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, server.router)
}
