package api

import (
	"net/http"

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

// ContextKey is used for go context keys
type ContextKey string

const (
	// IAMContextKey represents the context used for IAM data
	IAMContextKey ContextKey = "iam"
)

type httpHandler struct {
	t        *thrap.Thrap
	projects *thrap.Projects
}

// This is handled by the middleware
func (h *httpHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	val := ctx.Value(IAMContextKey)

	writeJSONResponse(w, val, nil)
}

func (h *httpHandler) handleOptionsLogin(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
}
