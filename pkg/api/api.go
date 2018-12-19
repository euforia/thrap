package api

import (
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/euforia/thrap/pkg/thrap"
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
	data, err := Asset("build/swagger.json")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
