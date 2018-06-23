package thrap

import (
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
)

type BasicStackConfig struct {
	Name      string
	Language  thrapb.LanguageID
	DataStore string
	WebServer string
}

// NewBasicStack builds a skeleton stack inferring as much information as
// possible to assemble the stack.  It only inits the object with the minimal
// valid defaults
func NewBasicStack(c *BasicStackConfig, pks *packs.Packs) (*thrapb.Stack, error) {

	stack := thrapb.Stack{
		ID:   c.Name,
		Name: c.Name,
	}

	comps := map[string]*thrapb.Component{
		consts.DefaultAPICompID: makeDevComp(c.Name, consts.DefaultAPICompID, c.Language),
	}

	if c.DataStore != "none" {
		//dsc := conf.DataStores[dataStore]
		dspacks := pks.Datastore()
		ds, err := dspacks.Load(c.DataStore)
		if err != nil {
			return nil, err
		}

		comp := thrapb.NewComponent(ds.Image, ds.DefaultVersion, thrapb.CompTypeDatastore)
		comps[consts.DefaultDSCompID] = comp
	}

	if c.WebServer != "none" {
		//wss := conf.WebServers[webServer]
		wpacks := pks.Web()
		ws, err := wpacks.Load(c.WebServer)
		if err != nil {
			return nil, err
		}
		comp := thrapb.NewComponent(ws.Image, ws.DefaultVersion, thrapb.CompTypeWeb)
		comps[consts.DefaultWebCompID] = comp
	}

	stack.Components = comps

	return &stack, nil
}

func makeDevComp(stackID, compID string, lang thrapb.LanguageID) *thrapb.Component {
	return &thrapb.Component{
		ID:       compID,
		Name:     "${" + vars.RegistryAddr + "}/" + stackID + "/" + compID,
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
			Format:      "hcl",
		},
	}
}
