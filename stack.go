package thrap

import (
	"github.com/euforia/thrap/config"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
)

// NewBasicStack builds a skeleton stack inferring as much information as
// possible to assemble the stack.  It only inits the object with the minimal
// valid defaults
func NewBasicStack(lang thrapb.LanguageID, name, dataStore, webServer string, conf *config.Config) *thrapb.Stack {

	stack := thrapb.Stack{
		ID:   name,
		Name: name,
	}

	comps := map[string]*thrapb.Component{
		consts.DefaultAPICompID: makeDevComp(name, consts.DefaultAPICompID, lang),
	}

	if dataStore != "none" {
		dsconf := conf.DataStores[dataStore]
		comps[consts.DefaultDSCompID] = makeComp(dsconf)
	}

	if webServer != "none" {
		wss := conf.WebServers[webServer]
		wsComp := makeComp(wss)
		comps[consts.DefaultWebCompID] = wsComp
	}

	stack.Components = comps

	return &stack
}

func makeComp(conf *config.ImageVersionConfig) *thrapb.Component {
	return &thrapb.Component{
		Name:    conf.Image,
		Type:    thrapb.CompTypeWeb,
		Version: conf.DefaultVersion,
	}
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
			// Context:    consts.DefaultBuildContext,
		},
		Env: &thrapb.Envionment{
			File: consts.DefaultEnvFile,
			Vars: map[string]string{
				consts.EnvVarVersion: "",
			},
		},
		Secrets: &thrapb.Secrets{
			Destination: consts.DefaultSecretsFile,
			Format:      "hcl",
		},
	}
}
