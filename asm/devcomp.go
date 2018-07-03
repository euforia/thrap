package asm

import (
	"path/filepath"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
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
	// language pack for component
	pack *packs.DevPack

	// devpack loaded scope vars
	vars scope.Variables

	// Component being assembled
	comp *thrapb.Component

	// Parsed docker file
	dockerfile    *dockerfile.Dockerfile
	dockerignores []string

	// Assembled files after rendering
	files map[string][]byte
}

// NewDevCompAsm returns a new initialized DevCompAsm
func NewDevCompAsm(c *thrapb.Component, langpack *packs.DevPack) *DevCompAsm {
	asm := &DevCompAsm{
		comp:          c,
		pack:          langpack,
		files:         make(map[string][]byte),
		dockerignores: make([]string, 0),
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

// Dockerfile returns the parsed dockerfile associated to the component
func (asm *DevCompAsm) Dockerfile() *dockerfile.Dockerfile {
	return asm.dockerfile
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

func (asm *DevCompAsm) addArgAndEnvToDockerfile() {
	cenv := asm.comp.Env
	ei := &dockerfile.Env{Vars: make(map[string]string, len(cenv.Vars))}
	for k := range cenv.Vars {
		ei.Vars[k] = "${" + k + "}"
	}
	asm.dockerfile.AddInstruction(0, ei)

	// Add all env vars as build args to build stage

	for k := range cenv.Vars {
		ai := &dockerfile.Arg{Name: k}
		asm.dockerfile.AddInstruction(0, ai)
	}

	// Add all env vars to the final image artifact
	pei := &dockerfile.Env{Vars: make(map[string]string, len(cenv.Vars))}
	for k := range cenv.Vars {
		pei.Vars[k] = ""
	}
	asm.dockerfile.AddInstruction(1, pei)
}

func (asm *DevCompAsm) addSecretsVolumeToDockerfile(i int, workdir *dockerfile.WorkDir) error {
	vol := &dockerfile.Volume{
		Paths: []string{filepath.Join(workdir.Path, asm.comp.Secrets.Destination)},
	}

	return asm.dockerfile.AddInstruction(i, vol)
}

func (asm *DevCompAsm) ensureDockerfileWorkdir(idx int) (*dockerfile.WorkDir, error) {
	stage := asm.dockerfile.Stages[idx]
	wd, _ := stage.GetOp(dockerfile.KeyWorkDir)

	if wd == nil {
		wd = &dockerfile.WorkDir{Path: consts.DefaultWorkDir}
		err := asm.dockerfile.AddInstruction(idx, wd)
		if err != nil {
			return nil, err
		}
	}

	return wd.(*dockerfile.WorkDir), nil
}

// func (asm *DevCompAsm) addSecretsVolumes(workdirs []*dockerfile.WorkDir) (err error) {
// 	for i, workdir := range workdirs {
//
// 		err = asm.addSecretsVolume(i, workdir)
// 		if err != nil {
// 			break
// 		}
// 	}
// 	return
// }

// addSecretVolToStage adds a volume instruction to the given stage.  It declares
// the volume relative to the WORKDIR specified. If WORKDIR is not specified it
// defaults to /
// func (asm *DevCompAsm) addSecretVolToStage(stageIndex int) error {
// 	workdir, err := asm.ensureWorkdir(stageIndex)
// 	if err != nil {
// 		return err
// 	}
//
// 	vol := &dockerfile.Volume{
// 		Paths: []string{filepath.Join(workdir.Path, asm.comp.Secrets.Destination)},
// 	}
//
// 	return asm.dockerfile.AddInstruction(stageIndex, vol)
// }
