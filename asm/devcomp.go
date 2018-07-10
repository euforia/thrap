package asm

import (
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/packs"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
	"github.com/hashicorp/hil/ast"
)

const defaultSecretsContent = `# HCL formatted file
# Examples:
#
# key = "value"
# foo.bar = 5
# feature.enabled = false
`

// DevCompAsm is dev component assembler i.e. buildable component with or without
// language specification
type DevCompAsm struct {
	*BuildCompAsm
	// language pack for component
	pack *packs.DevPack
	// Assembled files after rendering
	files map[string][]byte
}

// NewDevCompAsm returns a new initialized DevCompAsm
func NewDevCompAsm(c *thrapb.Component, langpack *packs.DevPack) *DevCompAsm {
	asm := &DevCompAsm{
		BuildCompAsm: &BuildCompAsm{
			comp:          c,
			dockerignores: make([]string, 0),
		},
		pack:  langpack,
		files: make(map[string][]byte),
	}

	asm.init()

	return asm
}

func (asm *DevCompAsm) init() {
	if asm.pack == nil {
		return
	}

	asm.vars = asm.pack.ScopeVars()

	// Default is set.  Only set if provided
	if asm.comp.Language.Version() != "" {

		asm.vars[vars.LangVersion] = ast.Variable{
			Value: asm.comp.Language.Version(),
			Type:  ast.TypeString,
		}
	}

}

func (asm *DevCompAsm) assemblePack(svars scope.Variables) (err error) {
	if asm.pack == nil {
		return
	}

	asm.dockerignores = append(asm.dockerignores, asm.pack.IgnoreFiles...)

	asm.dockerfile, err = asm.pack.Dockerfile(svars)
	if err == nil {
		asm.files, err = asm.pack.Files(svars)
	}

	return
}

// Assemble assembles all components assets, normalizing values using the given
// scope variables
func (asm *DevCompAsm) Assemble(variables scope.Variables) (err error) {
	scopeVars := make(scope.Variables, len(variables)+len(asm.vars))
	for k, v := range variables {
		scopeVars[k] = v
	}
	// These take precedence
	for k, v := range asm.vars {
		scopeVars[k] = v
	}

	if err = asm.assemblePack(scopeVars); err != nil {
		return err
	}

	// Dockerfile workdir check
	stages := asm.dockerfile.Stages
	lastIdx := len(stages) - 1

	workdirs := make([]*dockerfile.WorkDir, len(stages))
	for i := range stages {
		workdirs[i], err = asm.ensureDockerfileWorkdir(i)
		if err != nil {
			return err
		}
	}

	// Add all env vars as build args
	asm.addArgAndEnvToDockerfile()

	if asm.comp.HasSecrets() {
		// Local file
		asm.files[asm.comp.Secrets.Destination] = []byte(defaultSecretsContent)
		asm.dockerignores = append(asm.dockerignores, asm.comp.Secrets.Destination)
		err = asm.addSecretsVolumeToDockerfile(lastIdx, workdirs[lastIdx])
	}

	return err
}
