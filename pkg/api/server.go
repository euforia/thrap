package api

import (
	"log"
	"net"
	"net/http"

	"github.com/euforia/thrap/pkg/thrap"
	"github.com/gorilla/mux"
)

// Server is the REST api interface
type Server struct {
	// api request router
	router *mux.Router
	// api handler
	handler *httpHandler

	log *log.Logger
}

// NewServer returns a new API server
func NewServer(t *thrap.Thrap, logger *log.Logger) *Server {
	server := &Server{
		router: mux.NewRouter(),
		log:    logger,
		handler: &httpHandler{
			t:        t,
			projects: thrap.NewProjects(t),
			uiPrefix: "/ui",
		},
	}

	server.registerHandlers()
	if t.IAMEnabled() {
		server.enableAuthHandlers()
	} else {
		server.log.Println("IAM middleware DISABLED")
	}

	return server
}

// Serve starts serving the registered handlers on the given listener
func (server *Server) Serve(ln net.Listener) error {
	return http.Serve(ln, server.router)
}

func (server *Server) enableAuthHandlers() {
	server.router.HandleFunc("/v1/login", server.handler.handleOptionsLogin).Methods("OPTIONS")
	server.router.HandleFunc("/v1/login", server.handler.handleLogin).Methods("POST")
	server.router.Use(server.authMiddleware)
}

func (server *Server) registerHandlers() {
	server.router.PathPrefix("/ui/").HandlerFunc(server.handler.handleUI)
	// server.router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("/Users/abs/workbench/platform-ui/build"))))

	// No auth for profile operations
	server.router.HandleFunc("/v1/profiles", server.handler.handleListProfiles).Methods("GET")
	server.router.HandleFunc("/v1/profile/{id}", server.handler.handleProfile)

	//
	// Auth'd endpoints
	//
	server.registerHandler("/v1/projects", server.handler.handleListProjects)
	server.registerHandler("/v1/project/{id}", server.handler.handleProject)

	server.registerHandler("/v1/project/{pid}/deployments", server.handler.handleListDeployments)
	server.registerHandler("/v1/project/{pid}/deployment/spec", server.handler.handleDeploymentSpec)
	server.registerHandler("/v1/project/{pid}/deployment/{eid}/{iid}", server.handler.handleDeployment)
}

func (server *Server) registerHandler(path string, h http.HandlerFunc) {
	// server.router.HandleFunc(path, server.handleRequestAuth(h))
	server.router.HandleFunc(path, h)
}

func (server *Server) authMiddleware(next http.Handler) http.Handler {
	return &authHandler{
		next: next,
		t:    server.handler.t,
		log:  server.log,
	}
}