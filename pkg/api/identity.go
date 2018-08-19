package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/euforia/thrap/pkg/identity"
	"github.com/euforia/thrap/thrapb"

	"github.com/gorilla/mux"
)

type identities struct {
	store *identity.Identity
}

func (api *identities) list(w http.ResponseWriter, r *http.Request) {
	list := make([]*thrapb.Identity, 0)

	err := api.store.Iter("", func(ident *thrapb.Identity) error {
		list = append(list, ident)
		return nil
	})

	writeJSONResponse(w, list, err)
}

func (api *identities) identity(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var (
		resp *thrapb.Identity
		err  error
	)

	switch r.Method {
	case "GET":
		resp, err = api.get(w, r)

	case "POST":
		resp, err = api.register(w, r)

	case "DELETE":
		err = api.delete(w, r)

	default:
		w.WriteHeader(405)
		return
	}

	writeJSONResponse(w, resp, err)
}

func (api *identities) get(w http.ResponseWriter, r *http.Request) (*thrapb.Identity, error) {
	identID := mux.Vars(r)["id"]
	return api.store.Get(identID)
}

func (api *identities) register(w http.ResponseWriter, r *http.Request) (*thrapb.Identity, error) {
	identID := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var ident thrapb.Identity
	err = json.Unmarshal(body, &ident)
	if err != nil {
		return nil, err
	}
	ident.ID = identID

	id, _, err := api.store.Register(&ident)
	return id, err
}

func (api *identities) delete(w http.ResponseWriter, r *http.Request) error {
	identID := mux.Vars(r)["id"]
	return api.store.Delete(identID)
}
