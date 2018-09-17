package thrap

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/euforia/thrap/pkg/provider/orchestrator"
	"github.com/euforia/thrap/pkg/provider/registry"
	"github.com/euforia/thrap/pkg/provider/secrets"
	"github.com/euforia/thrap/thrapb"
)

var (
	errInvalidImageFormat = errors.New("invalid image format")
	errSecretsNotLoaded   = errors.New("secrets not loaded")
)

const (
	// SecretsPathVarName is the name of the secrets path variable made available
	// during deployment
	SecretsPathVarName = "SECRETS_PATH"
	// InstanceVarName is the deployment instance name
	InstanceVarName = "INSTANCE"
	// ProjectVarName is the project id made available during the deploy
	ProjectVarName = "PROJECT"
	// RegistryVarName is the registry name configured with the engine
	RegistryVarName = "REGISTRY"
	// DeployVerVarName is the deployment version meta variable. This is
	// incremented on each deploy
	DeployVerVarName = "DEPLOY_VERSION"
)

// Engine interface implements an engine to run deployments
type Engine interface {
	Profile() thrapb.Profile
	PrepareDeploy(*provider.Request) (orchestrator.PreparedDeployment, error)
	SetupSecrets(project string, data map[string]interface{}) error
	// Returns the secrets path given the project and instance name
	SecretsPath(project string, instance string) string
	SeedSecrets(*provider.Request) error
	Deploy(context.Context, orchestrator.PreparedDeployment, orchestrator.RequestOptions) error
}

type engine struct {
	// Profile used to load the engine
	profile *thrapb.Profile

	// Orchestrator based on profile
	o orchestrator.Orchestrator

	// Additional registry to the default docker hub registry which has
	// no registry name prefix
	r registry.Registry
	// Default for images with no registry
	dockerHub registry.Registry

	// Secrets backend
	sec secrets.Secrets

	log *log.Logger
}

func (eng *engine) Profile() thrapb.Profile {
	return *eng.profile.Clone()
}

// SetupSecrets populates the given key-value secrets to the template path
func (eng *engine) SetupSecrets(projID string, data map[string]interface{}) error {
	err := eng.sec.Create(projID)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		tmplKey := filepath.Join(projID, "template")
		return eng.sec.Set(tmplKey, data)
	}

	return nil
}

// returns absolute path
func (eng *engine) SecretsPath(projID, instID string) string {
	key := eng.secretsRelativePath(projID, instID)
	return eng.sec.SecretsPath(key)
}

// returns the relative path
func (eng *engine) secretsRelativePath(projID, instID string) string {
	return filepath.Join(projID, "instance", instID)
}

// SeedSecrets copies the secrets from the templates into the given instance
// in the request
func (eng *engine) SeedSecrets(req *provider.Request) error {
	projID := req.Project.ID
	instID := req.Deployment.Name

	tmplKey := filepath.Join(projID, "template")

	mm, err := eng.sec.RecursiveGet(tmplKey)
	if err != nil {
		return err
	}

	instKey := eng.secretsRelativePath(projID, instID)

	eng.log.Printf("Seeding secrets from=%s ---> to=%s", tmplKey, instKey)

	for k, v := range mm {
		if len(v) == 0 {
			continue
		}

		baseKey := strings.TrimPrefix(k, tmplKey)
		newKey := filepath.Join(instKey, baseKey)

		eng.log.Printf("Seeding secret=%s", newKey)
		er := eng.sec.Set(newKey, v)
		if er != nil {
			eng.log.Printf("Error seeding secret=%s: %v", newKey, er)
			err = er
		}
	}

	return err
}

func (eng *engine) PrepareDeploy(req *provider.Request) (orchestrator.PreparedDeployment, error) {
	if regName := eng.r.Name(); regName != "" {
		req.Deployment.Profile.Meta[RegistryVarName] = regName
	}

	prepared, err := eng.o.PrepareDeploy(req)
	if err != nil {
		return nil, err
	}

	err = eng.checkArtifacts(prepared)
	if err != nil {
		return prepared, err
	}

	err = eng.checkSecrets(req.Project.ID, req.Deployment.Name)

	//
	// All other checks go here
	//

	return prepared, err
}

func (eng *engine) Deploy(ctx context.Context,
	prepared orchestrator.PreparedDeployment, opts orchestrator.RequestOptions) error {

	_, err := eng.o.Deploy(ctx, prepared, opts)

	return err
}

func (eng *engine) checkSecrets(projID, instID string) error {
	path := eng.secretsRelativePath(projID, instID)
	eng.log.Printf("Checking secrets project=%s instance=%s path=%s",
		projID, instID, path)

	_, err := eng.sec.Get(path)

	return err
}

func (eng *engine) checkArtifacts(p orchestrator.PreparedDeployment) error {
	regName := eng.r.Name()
	arts := p.Artifacts()

	var err error
	for _, art := range arts {

		imgName := strings.TrimPrefix(art, regName+"/")
		parts := strings.Split(imgName, ":")
		lp := len(parts)

		switch lp {
		case 1:
			_, err = eng.getImageManifest(regName, imgName, "latest")
		case 2:
			_, err = eng.getImageManifest(regName, parts[0], parts[1])

		default:
			// return errors.New("invalid image format: " + art)
			err = errInvalidImageFormat
		}

		if err != nil {
			return errors.Wrap(err, art)
		}

	}
	return nil
}

func (eng *engine) getImageManifest(reg, img, tag string) (interface{}, error) {
	eng.log.Printf("Checking image=%s/%s:%s", reg, img, tag)
	if reg != "" {
		return eng.r.GetImageManifest(img, tag)
	}
	// Check docker hub as default
	return eng.dockerHub.GetImageManifest(img, tag)
}

func (eng *engine) validateResourceProviders() error {
	profile := eng.profile
	// if eng.sec == nil {
	// 	return fmt.Errorf("unsupported secrets: %s", profile.Secrets)
	// }
	if eng.o == nil {
		return fmt.Errorf("unsupported orchestrator: %s", profile.Orchestrator)
	}
	if eng.r == nil {
		return fmt.Errorf("unsupported registry: %s", profile.Registry)
	}
	return nil
}
