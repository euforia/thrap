package api

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/euforia/thrap/pkg/thrap"
)

const (
	// DescContentTypeMoldHCL is the legacy type to be deprecated
	DescContentTypeMoldHCL = "application/vnd.thrap.mold.deployment.descriptor.v1+hcl"
	// DescContentTypeNomadHCL is a nomad hcl file
	DescContentTypeNomadHCL = "application/vnd.thrap.nomad.deployment.descriptor.v1+hcl"
	// DescContentTypeNomadJSON is json object
	DescContentTypeNomadJSON = "application/vnd.thrap.nomad.deployment.descriptor.v1+json"
)

const (
	// TokenHeader is the auth token header key
	TokenHeader = "X-Vault-Token"
	// ProfileHeader is the profile header key
	ProfileHeader = "Thrap-Profile"
)

type httpHandler struct {
	t        *thrap.Thrap
	projects *thrap.Projects
	uiPrefix string
}

// This is handled by the middleware
func (h *httpHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := ctx.Value(thrap.IAMContextKey)

	writeJSONResponse(w, val, nil)
}

func (h *httpHandler) handleOptionsLogin(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.WriteHeader(200)
}

func (h *httpHandler) handleUI(w http.ResponseWriter, r *http.Request) {
	var (
		upath = strings.TrimPrefix(r.URL.Path, h.uiPrefix)
		fpath string
	)

	switch {
	case strings.HasPrefix(upath, "/static/"):
		fpath = filepath.Join("build", upath)
	default:
		fpath = filepath.Join("build", "index.html")
	}

	data, err := Asset(fpath)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}

	contentType := mime.TypeByExtension(path.Ext(fpath))
	w.Header().Add("Content-Type", contentType)
	w.WriteHeader(200)
	w.Write(data)
}

func (h *httpHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	w.WriteHeader(200)
}

func (h *httpHandler) handleSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	w.WriteHeader(200)
	w.Write([]byte(swaggerJSON))
}
