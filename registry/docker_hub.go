package registry

import (
	"errors"
	"fmt"

	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/libtrust"

	"github.com/euforia/docker-registry-client/registry"

	"github.com/euforia/thrap/config"
)

const defaultDockerRegAddr = "https://registry.hub.docker.com"

type dockerHub struct {
	id  string
	url string
	reg *registry.Registry
}

func (hub *dockerHub) ID() string {
	return hub.id
}

// Initialize the registry provider
func (hub *dockerHub) Init(rconf *config.RegistryConfig) error {
	hub.url = defaultDockerRegAddr

	conf := rconf.Config

	if val, ok := conf["url"]; ok {
		v, ok := val.(string)
		if !ok {
			return errors.New("registry url invalid")
		}
		hub.url = v
	}

	dockerHubConf := &registry.Config{}

	if val, ok := conf["user"]; ok {
		dockerHubConf.Username, ok = val.(string)
		if !ok {
			return errors.New("registry user invalid")
		}
	}

	if val, ok := conf["password"]; ok {
		dockerHubConf.Password, ok = val.(string)
		if !ok {
			return errors.New("registry password invalid")
		}
	}

	dreg, err := registry.New(hub.url, dockerHubConf)
	if err == nil {
		hub.reg = dreg
	}
	return err
}

// Create a new repository
func (hub *dockerHub) Create(name string) (interface{}, error) {

	mfest := schema1.Manifest{
		Versioned: manifest.Versioned{SchemaVersion: 2},
	}
	// smf, err := schema2.FromStruct(mfest)
	key, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return nil, err
	}

	signedManifest, err := schema1.Sign(&mfest, key)
	if err != nil {
		return nil, err
	}

	err = hub.reg.PutManifest(name, "latest", signedManifest)
	return signedManifest, err
}

func (hub *dockerHub) ImageName(name string) string {
	return name
}

// get repo info
func (hub *dockerHub) Get(name string) (interface{}, error) {
	return nil, fmt.Errorf("TBI")
}

// // Type of registry. container/deployment
// func (hub *dockerHub) Type() Type {
// 	return TypeContainer
// }

// Get a repository manifest
func (hub *dockerHub) GetManifest(name, tag string) (interface{}, error) {
	return hub.reg.ManifestV2(name, tag)
}
