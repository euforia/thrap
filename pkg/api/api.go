package api

import (
	"net"
	"net/http"

	"github.com/euforia/kvdb"

	"github.com/euforia/thrap/pkg/project"
	"github.com/gorilla/mux"
)

// Config is the api server config
type Config struct {
	Datastore kvdb.Datastore
}

type httpHandler struct {
	projs *project.Projects
}

// Server is the REST api interface
type Server struct {
	// api request router
	router *mux.Router
	// api handler
	handler *httpHandler
}

// NewServer returns a new API server
func NewServer(conf *Config) *Server {
	server := &Server{
		router: mux.NewRouter(),
		handler: &httpHandler{
			projs: project.NewProjects(conf.Datastore),
		},
		// ident:  &identities{},
	}

	// server.initHandlers(ds)
	server.registerHandlers()

	return server
}

// func (server *Server) initHandlers(ds *storage.Datastore) {
// server.proj.store, _ = ds.Project()

// conf := &identity.Config{}
// conf.Storage, _ = ds.Identity()
// server.ident.store = identity.New(conf)
// }

func (server *Server) registerHandlers() {
	server.router.HandleFunc("/projects", server.handler.handleListProjects)
	server.router.HandleFunc("/project/{id}", server.handler.handleProject)

	// server.router.HandleFunc("/identities", server.ident.list)
	// server.router.HandleFunc("/identity/{id}", server.ident.identity)

	server.router.HandleFunc("/project/{pid}/deployments", server.handler.handleListDeployments)

	server.router.HandleFunc("/project/{pid}/deployment/spec", server.handler.handleDeploymentSpec)
	server.router.HandleFunc("/project/{pid}/deployment/{eid}/{did}", server.handler.handleDeployment)

	// server.router.HandleFunc("/project/{pid}/deployment/{eid}", server.deploy.listEnvDeployments)
}

// Serve starts serving the registered handlers on the given listener
func (server *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, server.router)
}
