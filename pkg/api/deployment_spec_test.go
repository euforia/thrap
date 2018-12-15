package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/euforia/thrap/pkg/credentials"
	"github.com/euforia/thrap/pkg/thrap"
	"github.com/pkg/errors"
)

func newTestServer() (*Server, error) {
	conf := &thrap.Config{
		ConfigDir: "../../.thrap",
	}
	if err := conf.Validate(); err != nil {
		return nil, errors.Wrap(err, "new test server")
	}

	var err error
	credsFile := filepath.Join(conf.ConfigDir, "creds.hcl")
	conf.Credentials, err = credentials.ReadCredentials(credsFile)
	if err != nil {
		return nil, err
	}

	thp, err := thrap.New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "new test server")
	}

	return NewServer(thp, &log.Logger{}), nil
}

func executeRequest(srv *Server, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	srv.router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestListNoSpecs(t *testing.T) {
	srv, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "/v1/project/vfs/deployment/spec", nil)
	resp := executeRequest(srv, req)

	checkResponseCode(t, http.StatusNotFound, resp.Code)
}
