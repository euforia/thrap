package asm

import (
	"errors"
	"strings"

	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
)

var (
	errLanguageRequired = errors.New("language required")
)

// BasicStackConfig holds a configurations to build canned stacks.  Currently
// supports a dev, datastore and web component
type BasicStackConfig struct {
	Name       string
	Language   thrapb.LanguageID
	DataStores []string
	Proxy      string
}

// NewBasicStack builds a skeleton stack inferring as much information as
// possible to assemble the stack.  It only inits the object with the minimal
// valid defaults
func NewBasicStack(c *BasicStackConfig, pks *packs.Packs) (*thrapb.Stack, error) {
	if c.Language == "" {
		return nil, errLanguageRequired
	}

	stack := thrapb.Stack{
		ID:   c.Name,
		Name: c.Name,
	}

	devComp := makeDevComp(consts.DefaultAPICompID, c.Language)
	comps := map[string]*thrapb.Component{
		consts.DefaultAPICompID: devComp,
	}

	dspacks := pks.Datastore()
	for _, datastore := range c.DataStores {
		var dsComp *thrapb.Component
		ds, err := dspacks.Load(datastore)
		if err != nil {
			return nil, err
		}

		id := ds.Name

		dsComp = thrapb.NewComponent(ds.Image, ds.DefaultVersion, thrapb.CompTypeDatastore)
		dsComp.ID = id
		comps[id] = dsComp

		// Add env var for host, port and addr to dev component
		ev := defaultCompEnvVars(id)
		for k, v := range ev {
			devComp.Env.Vars[k] = v
		}
	}

	var wsComp *thrapb.Component
	if c.Proxy != "none" && c.Proxy != "" {
		wpacks := pks.Web()
		ws, err := wpacks.Load(c.Proxy)
		if err != nil {
			return nil, err
		}
		wsComp = thrapb.NewComponent(ws.Image, ws.DefaultVersion, thrapb.CompTypeWeb)
		comps[consts.DefaultWebCompID] = wsComp

		// Add env var for host, port and addr of dev to web component (the just created)
		wsComp.Env = &thrapb.Envionment{
			Vars: defaultCompEnvVars(consts.DefaultAPICompID),
		}
	}

	stack.Components = comps

	return &stack, nil
}

func defaultCompEnvVars(pre string) map[string]string {
	upre := strings.ToUpper(pre)
	return map[string]string{
		upre + "_CONTAINER_IP": "${" + consts.CompVarPrefixKey + "." + pre + ".container.ip}",
	}
}

func makeDevComp(compID string, lang thrapb.LanguageID) *thrapb.Component {
	return &thrapb.Component{
		ID: compID,
		// registry info and parent are dynamically prepended based on profile
		Name:     compID,
		Type:     thrapb.CompTypeAPI,
		Language: lang,
		Head:     true,
		Build: &thrapb.Build{
			Dockerfile: compID + "." + consts.DefaultDockerFile,
		},
		Env: &thrapb.Envionment{
			File: consts.DefaultEnvFile,
			Vars: map[string]string{
				consts.EnvVarVersion: "${" + vars.StackVersion + "}",
			},
		},
		Secrets: &thrapb.Secrets{
			Destination: consts.DefaultSecretsFile,
		},
	}
}
