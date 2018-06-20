package asm

import (
	"path/filepath"

	"github.com/euforia/pseudo/scope"
	"github.com/euforia/thrap/consts"
	"github.com/euforia/thrap/devpack"
	"github.com/euforia/thrap/dockerfile"
	"github.com/euforia/thrap/thrapb"
	"github.com/euforia/thrap/vars"
	"github.com/hashicorp/hil/ast"
)

type devCompAssembler struct {
	// devpack
	pack *devpack.DevPack
	// devpack loaded scope vars
	vars scope.Variables

	// Component being assembled
	comp *thrapb.Component

	// Parsed docker file
	dockerfile *dockerfile.Dockerfile
	// Assembled files after rendering
	files map[string][]byte
}

func newDevCompAssembler(c *thrapb.Component, langpack *devpack.DevPack) *devCompAssembler {
	asm := &devCompAssembler{
		comp:  c,
		pack:  langpack,
		files: make(map[string][]byte),
	}

	asm.init()

	return asm
}

func (asm *devCompAssembler) init() {
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

func (asm *devCompAssembler) assemblePack(svars scope.Variables) (err error) {
	if asm.pack == nil {
		return
	}

	asm.dockerfile, err = asm.pack.Dockerfile(svars)
	if err == nil {
		asm.files, err = asm.pack.Files(svars)
	}

	return
}

func (asm *devCompAssembler) assemble(variables scope.Variables) (err error) {
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

	stages := asm.dockerfile.Stages
	workdirs := make([]*dockerfile.WorkDir, len(stages))
	for i := range stages {
		workdirs[i], err = asm.ensureWorkdir(i)
		if err != nil {
			return err
		}
	}

	if asm.comp.HasSecrets() {
		// Local file
		asm.files[asm.comp.Secrets.Destination] = []byte{}
		err = asm.addSecretsVolumes(workdirs)
	}

	return err
}

func (asm *devCompAssembler) addSecretsVolumes(workdirs []*dockerfile.WorkDir) (err error) {
	for i, workdir := range workdirs {

		vol := &dockerfile.Volume{
			Paths: []string{filepath.Join(workdir.Path, asm.comp.Secrets.Destination)},
		}

		err = asm.dockerfile.AddInstruction(i, vol)
		if err != nil {
			break
		}
	}
	return
}

func (asm *devCompAssembler) ensureWorkdir(idx int) (*dockerfile.WorkDir, error) {
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

// addSecretVolToStage adds a volume instruction to the given stage.  It declares
// the volume relative to the WORKDIR specified. If WORKDIR is not specified it
// defaults to /
// func (asm *devCompAssembler) addSecretVolToStage(stageIndex int) error {
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
