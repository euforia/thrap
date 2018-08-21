package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"

	"github.com/euforia/thrap/crt"
	"github.com/euforia/thrap/metrics"
	"github.com/euforia/thrap/pkg/provider/registry"
	"github.com/euforia/thrap/thrapb"
)

// PublishOptions hold options to publishing an artifact
type PublishOptions struct {
	TagLatest bool
}

type artifactPublisher struct {
	reg registry.Registry
	crt *crt.Docker

	auth types.AuthConfig
}

func (pub *artifactPublisher) Publish(ctx context.Context, stack *thrapb.Stack,
	opts PublishOptions) (map[string]error, *metrics.Runtime, error) {

	runtime := (&metrics.Runtime{}).Start()
	defer runtime.End()

	// local publish
	var local bool

	err := pub.login(ctx)
	if err != nil {
		if err != registry.ErrAuthNotRequired {
			return nil, runtime, err
		}
		local = true
	}

	reqs := pub.buildPushRequests(stack, opts.TagLatest)
	resps := make(map[string]error, len(reqs))

	if local {
		// Nothing to do for local
		for image := range reqs {
			resps[pub.reg.ImageName(image)] = nil
		}
	} else {

		for image, req := range reqs {
			fmt.Printf("Publishing %s:\n\n", image)

			// Check repo exists
			_, err := pub.reg.GetRepo(req.Image)
			if err == nil {
				// Publish
				req.Image = pub.reg.ImageName(req.Image)
				err = pub.crt.ImagePush(ctx, req)
				fmt.Println()
			}
			resps[pub.reg.ImageName(image)] = err
		}
	}

	runtime.End()
	return resps, runtime, nil
}

func (pub *artifactPublisher) login(ctx context.Context) error {
	authConfig, err := pub.reg.GetAuthConfig()
	if err != nil {
		return err
	}

	pub.auth = authConfig
	if pub.auth.Auth != "" {
		err = pub.crt.RegistryLogin(ctx, types.AuthConfig{
			Auth: pub.auth.Auth,
		})
	} else {
		err = pub.crt.RegistryLogin(ctx, types.AuthConfig{
			Username:      pub.auth.Username,
			Password:      pub.auth.Password,
			ServerAddress: pub.auth.ServerAddress,
		})
	}

	if err != nil {
		return err
	}

	return err
}

func (pub *artifactPublisher) getRegistryAuth() string {
	b, _ := json.Marshal(&types.AuthConfig{
		Username: pub.auth.Username,
		Password: pub.auth.Password,
	})
	return base64.URLEncoding.EncodeToString(b)
}

func (pub *artifactPublisher) buildPushRequests(stack *thrapb.Stack, tagLatest bool) map[string]*crt.PushRequest {
	reqs := make(map[string]*crt.PushRequest, len(stack.Components)*2)
	for id, comp := range stack.Components {
		if !comp.IsBuildable() {
			continue
		}

		name := stack.ArtifactName(id)

		if tagLatest {
			reqs[name+":latest"] = &crt.PushRequest{
				Image:  name + ":latest",
				Output: os.Stdout,
				Options: types.ImagePushOptions{
					RegistryAuth: pub.getRegistryAuth(),
				},
			}
		}
		reqs[name+":"+comp.Version] = &crt.PushRequest{
			Image:  name,
			Tag:    comp.Version,
			Output: os.Stdout,
			Options: types.ImagePushOptions{
				RegistryAuth: pub.getRegistryAuth(),
			},
		}
	}

	return reqs
}
