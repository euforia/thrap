package asm

import (
	"io/ioutil"
	"path/filepath"

	"github.com/euforia/pseudo"
	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
)

// ComponentAssembler implements a component assembly interface
type ComponentAssembler interface {
	Dockerfile() *dockerfile.Dockerfile
	Assemble(scope.Variables) error
}

// BuildCompAsm is a component assembler with dockerfiles
type BuildCompAsm struct {
	//
	vars scope.Variables

	// Component being assembled
	comp *thrapb.Component

	// Parsed docker file
	dockerfile    *dockerfile.Dockerfile
	dockerignores []string

	// Working/Context directory
	ctxDir string

	// HICL language vm
	vm *pseudo.VM
}

func NewBuildCompAsm(ctxDir string, comp *thrapb.Component) *BuildCompAsm {
	return &BuildCompAsm{
		comp:   comp,
		vm:     pseudo.NewVM(),
		ctxDir: ctxDir,
	}
}

func (asm *BuildCompAsm) parseEvalBytes(d []byte, svars scope.Variables) (string, error) {

	result, err := asm.vm.ParseEval(string(d), svars)
	if err == nil {
		return result.Value.(string), nil
	}

	return "", err
}

// Dockerfile returns the parsed dockerfile associated to the component
func (asm *BuildCompAsm) Dockerfile() *dockerfile.Dockerfile {
	return asm.dockerfile
}

// Assemble assembles the component applying the scoped variables.  This is
// used when a component does not have a language pack
func (asm *BuildCompAsm) Assemble(variables scope.Variables) (err error) {
	scopeVars := make(scope.Variables, len(variables)+len(asm.vars))
	// Add supplied vars
	for k, v := range variables {
		scopeVars[k] = v
	}
	// Internal vars take precedence
	for k, v := range asm.vars {
		scopeVars[k] = v
	}

	// Assemble dockerfile
	dfile := filepath.Join(asm.ctxDir, asm.comp.Build.Dockerfile)
	err = asm.loadDockerfile(dfile)
	if err != nil {
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
		// asm.files[asm.comp.Secrets.Destination] = []byte(defaultSecretsContent)
		asm.dockerignores = append(asm.dockerignores, asm.comp.Secrets.Destination)
		err = asm.addSecretsVolumeToDockerfile(lastIdx, workdirs[lastIdx])
	}

	return err
}

// loadDockerfile loads and parses the dockerfile at the given path
func (asm *BuildCompAsm) loadDockerfile(fpath string) error {
	// fpath := filepath.Join(ctxDir, asm.comp.Build.Dockerfile)
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	raw, err := dockerfile.ParseBytes(b)
	if err != nil {
		return err
	}

	asm.dockerfile = dockerfile.ParseRaw(raw)
	return nil
}

func (asm *BuildCompAsm) addArgAndEnvToDockerfile() {
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

func (asm *BuildCompAsm) addSecretsVolumeToDockerfile(i int, workdir *dockerfile.WorkDir) error {
	vol := &dockerfile.Volume{
		Paths: []string{filepath.Join(workdir.Path, asm.comp.Secrets.Destination)},
	}

	return asm.dockerfile.AddInstruction(i, vol)
}

func (asm *BuildCompAsm) ensureDockerfileWorkdir(idx int) (*dockerfile.WorkDir, error) {
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
