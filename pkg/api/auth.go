package api

import (
	"context"
	"log"
	"net/http"

	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/thrap"
)

type authHandler struct {
	t    *thrap.Thrap
	next http.Handler
	log  *log.Logger
}

func (a *authHandler) authWriteRequest(w http.ResponseWriter, r *http.Request) {
	profID := r.Header.Get("Thrap-Profile")
	if profID == "" {
		setAccessControlHeaders(w)
		w.WriteHeader(401)
		w.Write([]byte("profile id required"))
		return
	}

	token := r.Header.Get("X-Vault-Token")
	if token == "" {
		setAccessControlHeaders(w)
		w.WriteHeader(401)
		return
	}

	resp, err := a.t.Authenticate(profID, token)
	if err != nil {
		setAccessControlHeaders(w)
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	creds := &credentials.Credentials{
		Secrets: map[string]map[string]string{
			profID: map[string]string{
				"token": token,
			},
		},
	}

	ctx := context.WithValue(r.Context(), thrap.IAMContextKey, resp)
	ctx = context.WithValue(ctx, thrap.CredsContextKey, creds)

	a.next.ServeHTTP(w, r.WithContext(ctx))
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recW := &responseRecorder{w, 0}

	switch r.Method {
	case "PUT", "POST", "PATCH", "DELETE":
		a.authWriteRequest(recW, r)

	default:
		a.next.ServeHTTP(recW, r)
	}

	a.log.Printf("%d %s %s", recW.status, r.Method, r.URL.Path)
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (w *responseRecorder) WriteHeader(c int) {
	w.status = c
	w.ResponseWriter.WriteHeader(c)
}
