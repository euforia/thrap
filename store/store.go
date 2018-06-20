package store

import (
	"errors"
	"time"

	"github.com/euforia/thrap/thrapb"
)

var (
	errIDMissing = errors.New("ID missing")
)

//
// type Config struct {
// 	ObjectStorage ObjectStorage
// 	HashFunc      func() hash.Hash
// }
//
// func DefaultConfig() *Config {
// 	return &Config{
// 		HashFunc: sha256.New,
// 	}
// }

// ObjectStorage implements a namespaced object storage interface
type ObjectStorage interface {
	CreateRef(namespace, ref string) ([]byte, *thrapb.ChainHeader, error)
	SetRef(namespace, ref string, robj *thrapb.ChainHeader) ([]byte, error)
	GetRef(namespace, ref string) (*thrapb.ChainHeader, []byte, error)
	Set(namespace string, obj Object) ([]byte, error)
	Get(namepace string, digest []byte, obj Object) error
}

type IdentityStore struct {
	st ObjectStorage
	//hf func() hash.Hash
}

func NewIdentityStore(objs ObjectStorage) *IdentityStore {
	return &IdentityStore{st: objs}
}

func (store *IdentityStore) Get(id string) (*thrapb.Identity, *thrapb.ChainHeader, error) {
	ref, _, err := store.st.GetRef(id, "latest")
	if err != nil {
		return nil, nil, err
	}

	var ident thrapb.Identity
	err = store.st.Get(id, ref.DataDigest, &ident)

	return &ident, ref, err
}

func (store *IdentityStore) Create(ident *thrapb.Identity) (*thrapb.Identity, *thrapb.ChainHeader, error) {
	if ident.ID == "" {
		return nil, nil, errIDMissing
	}

	prev, _, err := store.st.CreateRef(ident.ID, "latest")
	if err != nil {
		return nil, nil, err
	}

	header, err := store.setIdent(ident, "latest", prev)

	return ident, header, err
}

func (store *IdentityStore) Update(ident *thrapb.Identity) (*thrapb.Identity, *thrapb.ChainHeader, error) {
	_, prev, err := store.st.GetRef(ident.ID, "latest")
	if err != nil {
		return nil, nil, err
	}

	header, err := store.setIdent(ident, "latest", prev)
	return ident, header, err
}

//
// func (store *IdentityStore) Delete(id string) error {
//
// }

// func BuildDockerfile(stackID string, comp *thrapb.Component, lang *languages.Language) (*dockerfile.Dockerfile, error) {
// 	var (
// 		stages = make([]dockerfile.Stage, 2)
// 		err    error
// 	)
//
// 	// Working directory
// 	workdir := "/src"
// 	switch lang.Name {
// 	case "go":
// 		workdir = languages.DefaultGoWorkDir("github.com", "owner", stackID)
//
// 	}
//
// 	bid := comp.ID + ".build"
//
// 	stages[0] = dockerfile.Stage{
// 		// Language pack
// 		&dockerfile.From{
// 			Image: lang.GetDevImages()[0] + ":" + comp.Language.Version(),
// 			As:    bid,
// 		},
// 		nil,
// 		&dockerfile.WorkDir{Path: workdir},
// 		&dockerfile.Env{
// 			Vars: map[string]string{"APP_VERSION": ""},
// 		},
// 		&dockerfile.Expose{},
// 		// Copy context into container i.e source code
// 		&dockerfile.Copy{Source: ".", Destination: "."},
// 		&dockerfile.Comment{},
// 		&dockerfile.Comment{Text: "ADD ADDITIONAL INSTRUCTIONS HERE"},
// 		&dockerfile.Comment{},
// 		&dockerfile.Run{Command: "<BUILD COMMANDS>"},
// 	}
//
// 	pubWorkDir := &dockerfile.WorkDir{Path: "/"}
// 	stages[1] = dockerfile.Stage{
// 		&dockerfile.From{Image: lang.GetPubImages()[0], As: comp.ID},
// 		nil,
// 		pubWorkDir,
// 		&dockerfile.Env{
// 			Vars: map[string]string{"APP_VERSION": ""},
// 		},
// 		&dockerfile.Expose{},
// 		&dockerfile.Comment{},
// 		&dockerfile.Comment{Text: "ADD ADDITIONAL INSTRUCTIONS HERE"},
// 		&dockerfile.Comment{},
// 		&dockerfile.Run{Command: "<COMMANDS>"},
// 		&dockerfile.Copy{
// 			Options:     []string{"--from=" + bid},
// 			Source:      workdir + "/<FILES FROM FIRST STAGE>",
// 			Destination: "<DESTINATION IN THIS CONTAINER>",
// 		},
// 		&dockerfile.Cmd{Command: "<APP EXEC COMMAND", Args: []string{"AND", "ARGS>"}},
// 		&dockerfile.EntryPoint{},
// 	}
//
// 	if comp.HasSecrets() {
// 		sdir := comp.Secrets.Destination
// 		if !filepath.IsAbs(sdir) {
// 			sdir = filepath.Join(pubWorkDir.Path, comp.Secrets.Destination)
// 		}
//
// 		vol := &dockerfile.Volume{
// 			Paths: []string{sdir},
// 		}
// 		stages[0][1] = vol
// 		stages[1][1] = vol
// 	}
//
// 	return &dockerfile.Dockerfile{Stages: stages}, err
// }

func (store *IdentityStore) setIdent(ident *thrapb.Identity, ref string, prev []byte) (*thrapb.ChainHeader, error) {
	data, err := store.st.Set(ident.ID, ident)
	if err != nil {
		return nil, err
	}

	refobj := &thrapb.ChainHeader{
		Previous:   prev,
		DataDigest: data,
		Timestamp:  time.Now().UnixNano(),
	}

	_, err = store.st.SetRef(ident.ID, ref, refobj)
	return refobj, err
}
