package registry

import (
	"errors"

	"github.com/heroku/docker-registry-client/registry"
)

const defaultDockerRegAddr = "https://registry.hub.docker.com"

type dockerRegistry struct {
	url string
	reg *registry.Registry
}

func (hub *dockerRegistry) ID() string {
	return "docker"
}

// Initialize the registry provider
func (hub *dockerRegistry) Init(config Config) error {
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
func (hub *dockerRegistry) Type() Type {
	return TypeContainer
}

// Create a new repository
func (hub *dockerRegistry) Create(string) (interface{}, error) {
	return nil, errors.New("to be implemented")
}

// Get a repository manifest
func (hub *dockerRegistry) GetManifest(name, tag string) (interface{}, error) {
	return hub.reg.ManifestV2(name, tag)
}
