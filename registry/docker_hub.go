package registry

import (
	"errors"

	"github.com/docker/distribution/manifest"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/libtrust"
	"github.com/heroku/docker-registry-client/registry"
)

type dockerHub struct {
	url string
	reg *registry.Registry
}

func (hub *dockerHub) ID() string {
	return "docker"
}

// Initialize the registry provider
func (hub *dockerHub) Init(config Config) error {
	hub.url = defaultDockerRegAddr

	conf := config.Conf

	if val, ok := conf["url"]; ok {
		v, ok := val.(string)
		if !ok {
			return errors.New("registry url invalid")
		}
		hub.url = v
	}

	var user string
	if val, ok := conf["user"]; ok {
		user, ok = val.(string)
		if !ok {
			return errors.New("registry user invalid")
		}
	}

	var passwd string
	if val, ok := conf["password"]; ok {
		passwd, ok = val.(string)
		if !ok {
			return errors.New("registry password invalid")
		}
	}

	dreg, err := registry.New(hub.url, user, passwd)
	if err == nil {
		hub.reg = dreg
	}

	return err
}

// Type of registry. container/deployment
func (hub *dockerHub) Type() Type {
	return TypeContainer
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

// Get a repository manifest
func (hub *dockerHub) GetManifest(name, tag string) (interface{}, error) {
	return hub.reg.ManifestV2(name, tag)
}
